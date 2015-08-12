package main

import (
	"io"
	"sort"
	"crypto/sha1"
	"strings"
	"fmt"
	"net/http"
)

type Request struct {
	Token string
}

func (this *Request) checkSignature(r *http.Request) bool {
	timestamp := r.FormValue("timestamp")
	nonce := r.FormValue("nonce")
	sl := []string{this.Token,timestamp,nonce}
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
