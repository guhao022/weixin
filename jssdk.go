package weixin

import (
	"strings"
	"io/ioutil"
	"os"
	"path"
	"errors"
	"time"
	"sort"
	"strconv"
	"fmt"
	"crypto/sha1"
	"io"
	"log"
)

type JsapiTicket struct {
	AccessToken AccessToken
	TmpName	string
	LckName string
}

func (this *JsapiTicket) Fresh() (string, error) {
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
}
