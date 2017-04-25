package weixin

import (
	"crypto/sha1"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
)

type Request struct {
	Token string
	// request common fields
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	// message request fields
	Content      string
	MsgId        int64
	PicUrl       string
	MediaId      string
	Format       string
	ThumbMediaId string
	LocationX    float64 `xml:"Location_X"`
	LocationY    float64 `xml:"Location_Y"`
	Scale        float64
	Label        string
	Title        string
	Description  string
	Url          string
	Recognition  string
	// event request fields
	Event     string
	EventKey  string
	Ticket    string
	Latitude  float64
	Longitude float64
	Precision float64
}

func (this *Request) checkSignature(r *http.Request) bool {
	timestamp := r.FormValue("timestamp")
	nonce := r.FormValue("nonce")
	sl := []string{this.Token, timestamp, nonce}
	sort.Strings(sl)

	s := sha1.New()
	io.WriteString(s, strings.Join(sl, ""))
	sign := fmt.Sprintf("%x", s.Sum(nil))

	signature := r.FormValue("signature")

	if sign != signature {
		return false
	} else {
		return true
	}
}

func (this *Request) Valid(w http.ResponseWriter, r *http.Request) bool {
	r.ParseForm()
	if !this.checkSignature(r) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
		return false
	}
	if r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(r.FormValue("echostr")))
		return false
	}
	return true
}
