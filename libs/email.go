package libs

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/itering/subscan/util"
)

func SubscribeToAddressBook(email string) {
	config := map[string]string{"signType": "md5"}
	request := map[string]string{"address": email}
	config["appid"] = util.GetEnv("MAIL_APPID", "14637")
	config["appkey"] = util.GetEnv("MAIL_APPKEY", "")
	subscribeRun(request, config)
}

func SendToSubscribe(to string) {
	config := map[string]string{"signType": "md5"}
	request := map[string]string{"to": to, "project": "CMSw01"}
	config["appid"] = util.GetEnv("MAIL_APPID", "14637")
	config["appkey"] = util.GetEnv("MAIL_APPKEY", "")
	emailXSendRun(request, config)
}

func emailXSendRun(request map[string]string, config map[string]string) string {
	request["appid"] = config["appid"]
	request["timestamp"] = getTimeStamp()
	request["signature"] = createSignature(request, config)
	return httpPost("https://api.mysubmail.com/mail/xsend", request)
}

func subscribeRun(request map[string]string, config map[string]string) string {
	request["appid"] = config["appid"]
	request["timestamp"] = getTimeStamp()
	request["signature"] = createSignature(request, config)
	return httpPost("https://api.mysubmail.com/addressbook/mail/subscribe", request)
}

func httpGet(queryUrl string) string {
	u, _ := url.Parse(queryUrl)
	retStr, err := http.Get(u.String())
	if err != nil {
		return err.Error()
	}
	result, err := ioutil.ReadAll(retStr.Body)
	retStr.Body.Close()
	if err != nil {
		return err.Error()
	}
	return string(result)
}

func httpPost(queryUrl string, postData map[string]string) string {
	data, err := json.Marshal(postData)
	if err != nil {
		return err.Error()
	}
	body := bytes.NewBuffer([]byte(data))
	retStr, err := http.Post(queryUrl, "application/json;charset=utf-8", body)
	if err != nil {
		return err.Error()
	}
	result, err := ioutil.ReadAll(retStr.Body)
	retStr.Body.Close()
	if err != nil {
		return err.Error()
	}
	return string(result)
}

func getTimeStamp() string {
	resp := httpGet("https://api.submail.cn/service/timestamp.json")
	var dict map[string]interface{}
	err := json.Unmarshal([]byte(resp), &dict)
	if err != nil {
		return err.Error()
	}
	return strconv.Itoa(int(dict["timestamp"].(float64)))
}

func createSignature(request map[string]string, config map[string]string) string {
	appKey := config["appkey"]
	appId := config["appid"]
	signType := config["signType"]
	request["sign_type"] = signType
	keys := make([]string, 0, 32)
	for key := range request {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	strList := make([]string, 0, 32)
	for _, key := range keys {
		strList = append(strList, fmt.Sprintf("%s=%s", key, request[key]))
	}
	sigStr := fmt.Sprintf("%s%s%s%s%s", appId, appKey, strings.Join(strList, "&"), appId, appKey)
	if signType == "normal" {
		return appKey
	} else if signType == "md5" {
		myMd5 := md5.New()
		_, _ = io.WriteString(myMd5, sigStr)
		return fmt.Sprintf("%x", myMd5.Sum(nil))
	} else {
		h := sha1.New()
		_, _ = io.WriteString(h, sigStr)
		return fmt.Sprintf("%x", h.Sum(nil))
	}
}
