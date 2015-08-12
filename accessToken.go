package main

import (
	"strings"
	"io/ioutil"
	"os"
	"path"
	"errors"
	"time"
)

type AccessToken struct {
	AppId     string
	AppSecret string
	TmpName   string
	LckName   string
}

func (this *AccessToken) Fresh() (string, error){
	if this.TmpName == "" {
		this.TmpName = "accesstoken.tmp"
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

func (this *AccessToken) fetchAndStore() (string, error) {
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

func (this *AccessToken) storage(token string) error {
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

func (this *AccessToken) fetch() (string, error) {
	requestUrl := strings.Join([]string{Url, "token?grant_type=client_credential&appid=", this.AppId, "&secret=", this.AppSecret}, "")
	resp, err := Get(requestUrl)
	if err != nil {
		return "", err
	}
	return resp.AccessToken, nil
}

func (this *AccessToken) unlock() error {
	return os.Remove(this.LckName)
}

func (this *AccessToken) lock() error {
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

func (this *AccessToken) locked() bool {
	_, err := os.Stat(this.LckName)
	return !os.IsNotExist(err)
}

