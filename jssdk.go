package weixin

import (
	"strings"
	"errors"
	"time"
	"sync"
    "encoding/json"
    "net/http"
    "strconv"
    "log"
    "sort"
    "crypto/sha1"
    "fmt"
    "io"
)

type JsapiTicket struct {
    accessToken AccessToken
	ticket Response
	locker sync.RWMutex
}

func (this *JsapiTicket) fetch() string {
    this.locker.RLock()
    defer this.locker.RUnlock()

    return this.ticket.Ticket
}

func (this *JsapiTicket) getJsApiTicket() error {
	response := struct {
		Code      int    `json:"errcode"`
		Msg       string `json:"errmsg"`
		Ticket    string `json:"ticket"`
		ExpiresIn int64  `json:"expires_in"`
	}{}
    access_token, err := this.accessToken.Fresh()
    if err != nil {
        return err
    }

    requestUrl := strings.Join([]string{Url, "ticket/getticket?access_token=", access_token, "&type=jsapi"}, "")
    resp, err := http.Get(requestUrl)
    if err != nil {
        return err
    }

	json.NewDecoder(resp.Body).Decode(&response)
	if response.Code != 0 {
		return errors.New(response.Msg)
	}

	this.locker.Lock()
	defer this.locker.Unlock()
    this.ticket.ExpiresIn = response.ExpiresIn
    this.ticket.Ticket = response.Ticket
	return nil
}

func (this *JsapiTicket) Get() error {
    if err := this.getJsApiTicket(); err != nil {
        return err
    }

    return nil
}

func (this *JsapiTicket) Refresh(refresh bool) {
    //先获得一次,获得失败panic
    if err := this.Get(); err != nil {
        panic(err)
    }

    //如果开启了自动刷新，就自动刷新
    if refresh {
        go func() {
            for {
                if err := this.Get(); err != nil {
                    continue
                }
                time.Sleep(time.Second * time.Duration(this.ticket.ExpiresIn))
                //time.Sleep(1 * time.Second)
                log.Println("Ticket: ", this.fetch())
            }
        }()
    }
}

type JsSign struct {
    AppID	  string `json:"appID"`
    NonceStr  string `json:"nonceStr"`
    Timetamp  int64  `json:"timestamp"`
    Url       string `json:"url"`
    Signature string `json:"signature"`
}

func GetJsSign(url string) *JsSign{
    var jsTick JsapiTicket
    jsTick.Refresh(true)
    timestamp := time.Now().Unix()
    noncestr := string(RandomCreateBytes(16))
    jsapi_ticket := jsTick.fetch()

    sl := []string{"timestamp=" + strconv.Itoa(int(timestamp)), "noncestr=" + noncestr,"jsapi_ticket=" + jsapi_ticket, "url=" + url}
    sort.Strings(sl)
    sortStr := strings.Join(sl, "&")

    log.Println("sortStr: ", sortStr)

    s := sha1.New()
    io.WriteString(s, sortStr)
    sign := fmt.Sprintf("%x", s.Sum(nil))
    return &JsSign{jsTick.accessToken.AppId,noncestr, timestamp, url, sign}
}

/*func (this *JsapiTicket) Fresh() (string, error) {
	if this.TmpName == "" {
		this.TmpName = "ticket.tmp"
	}
	if this.LckName == "" {
		this.LckName = this.TmpName + ".lck"
	}
	for {
		if this.locked() {
			time.Sleep(time.Second)
			continue
		}
		break
	}
	fi, err := os.Stat(this.TmpName)
	if err != nil && !os.IsExist(err) {
		return this.fetchAndStore()
	}
	expires := fi.ModTime().Add(2 * time.Hour).Unix()
	if expires <= time.Now().Unix() {
		return this.fetchAndStore()
	}
	tmp, err := os.Open(this.TmpName)
	if err != nil {
		return "", err
	}
	defer tmp.Close()
	data, err := ioutil.ReadAll(tmp)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (this *JsapiTicket) fetchAndStore() (string, error) {
	if err := this.lock(); err != nil {
		return "", err
	}
	defer this.unlock()
	token, err := this.fetch()
	if err != nil {
		return "", err
	}
	if err := this.storage(token); err != nil {
		return "", err
	}
	return token, nil
}

func (this *JsapiTicket) storage(token string) error {
	path := path.Dir(this.TmpName)
	fi, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	}
	if !fi.IsDir() {
		return errors.New("path is not a directory")
	}
	tmp, err := os.OpenFile(this.TmpName, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer tmp.Close()
	if _, err := tmp.Write([]byte(token)); err != nil {
		return err
	}
	return nil
}

func (this *JsapiTicket) fetch() (string, error) {
	access_token, err := this.AccessToken.Fresh()
	if err != nil {
		return "", err
	}
	requestUrl := strings.Join([]string{Url, "ticket/getticket?access_token=", access_token, "&type=jsapi"}, "")
	resp, err := Get(requestUrl)
	if err != nil {
		return "", err
	}
	return resp.Ticket, nil
}

func (this *JsapiTicket) unlock() error {
	return os.Remove(this.LckName)
}

func (this *JsapiTicket) lock() error {
	path := path.Dir(this.LckName)
	fi, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	}
	if !fi.IsDir() {
		return errors.New("path is not a directory")
	}
	lck, err := os.Create(this.LckName)
	if err != nil {
		return err
	}
	lck.Close()
	return nil
}

func (this *JsapiTicket) locked() bool {
	_, err := os.Stat(this.LckName)
	return !os.IsNotExist(err)
}

type JsSign struct {
	AppID	  string `json:"appID"`
	NonceStr  string `json:"nonceStr"`
	Timetamp  int64  `json:"timestamp"`
	Url       string `json:"url"`
	Signature string `json:"signature"`
}

func GetJsSign(url string) *JsSign{
	var jsTick JsapiTicket
	timestamp := time.Now().Unix()
	noncestr := string(RandomCreateBytes(16))
	jsapi_ticket, err := jsTick.Fresh()
	log.Println("Ticket: ", jsapi_ticket)
	if err != nil {
		panic(err)
	}
	sl := []string{"timestamp="+strconv.Itoa(int(timestamp)),"noncestr="+noncestr,"jsapi_ticket="+jsapi_ticket,"url="+url}
	sort.Strings(sl)
	sortStr := strings.Join(sl, "&")

	s := sha1.New()
	io.WriteString(s, sortStr)
	sign := fmt.Sprintf("%x", s.Sum(nil))
	return &JsSign{jsTick.AccessToken.AppId,noncestr, timestamp, url, sign}
}*/
