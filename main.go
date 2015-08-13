package main

import (
	"log"
	"net/http"
	"fmt"
)

const (
	AppID = "wxce6bcf1fd09e7c3c"
	appsecret = "6f14e784c4e49f5504ebb8d4802c5b61"
	token = "wordhi"
)

func main() {

	http.HandleFunc("/ticket", ticket)
	log.Println("监听端口 :9900...")
	err := http.ListenAndServe(":9900", nil)
	if err != nil {
		log.Println(err)
	}
}

func ticket(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	url := r.FormValue("url")
	if(url == ""){
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(http.StatusText(http.StatusAccepted)))
	}

	mp := New(token, appID, appsecret)

	t, _ := mp.AccessToken.Fresh()
	log.Println("AccessToken:",t)
	js_sign := GetJsSign(url)
	log.Println("JsSign: ", js_sign.Signature)
	log.Println("JsSignUrl: ", js_sign.Url)
	
	fmt.Println("====================================================================")
	OutputJson(w, js_sign)
}




/*func proc(w http.ResponseWriter, r *http.Request) {
	mp := New(token, appID, appsecret)
	f := mp.Request.Valid(w, r)
	log.Println(f)
}*/
