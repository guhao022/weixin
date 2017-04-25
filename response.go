package weixin

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Response struct {
	// error
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	// token
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`

	// media fields
	Type      string `json:"type"`
	MediaId   string `json:"media_id"`
	CreatedAt int64  `json:"created_at"`
	// ticket fields
	Ticket        string `json:"ticket"`
	ExpireSeconds int64  `json:"expire_seconds"`
}

func post(url string, bodyType string, body *bytes.Buffer) (*Response, error) {
	resp, err := http.Post(url, bodyType, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var rtn Response
	if err := json.Unmarshal(data, &rtn); err != nil {
		return nil, err
	}
	if rtn.ErrCode != 0 {
		return nil, errors.New(fmt.Sprintf("%d %s", rtn.ErrCode, rtn.ErrMsg))
	}
	return &rtn, nil
}

func get(url string) (*Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var rtn Response
	if err := json.Unmarshal(data, &rtn); err != nil {
		return nil, err
	}
	if rtn.ErrCode != 0 {
		return nil, errors.New(fmt.Sprintf("%d %s", rtn.ErrCode, rtn.ErrMsg))
	}
	return &rtn, nil
}
