package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	. "github.com/zt3862266/go/log"
)

type CallbackResult struct {
	Status int    `json:"error"`
	Msg    string `json:"msg"`
}

const CallbackResSuc = 0
const CallbackResFail = 1

func NotifyMsg(client *http.Client, url string, body []byte) int {

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		Warn("notifyMsg create request failed:%s", err)
		return CallbackResFail
	}
	req.Header.Set("Content-Type", "application/json")
	response, err := client.Do(req)

	if err != nil {
		Warn("request failed,url:%s ,err:%s", url, err)
		return CallbackResFail
	}
	defer response.Body.Close()
	statusCode := response.StatusCode
	retStr, err := ioutil.ReadAll(response.Body)
	Info("send post,url:%s,msg:%s,ret:%s", url, string(body), retStr)
	if err != nil {
		Warn("get response failed,url:%s,err:%s,statusCode:%d", url, err, statusCode)
		return CallbackResFail
	}
	var ret CallbackResult
	err = json.Unmarshal(retStr, &ret)
	if err != nil {
		Warn("json unmarshal failed,url:%s,err:%s,ret:%s", url, err, retStr)
		return CallbackResFail
	}
	if statusCode == 200 && ret.Status == CallbackResSuc {
		return CallbackResSuc
	}
	return CallbackResFail

}
