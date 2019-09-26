package dingding

import (
	"encoding/json"
	"errors"
	"strings"
	"subscan-end/utiles"
	"subscan-end/utiles/pusher"
)

const AccessUrl = "https://oapi.dingtalk.com/robot/send?access_token="

var DingClient = New()

type (
	Ding struct{}
	msg  struct {
		MsgType string `json:"msgtype"`
		Text    struct {
			Content string `json:"content"`
		} `json:"text"`
	}
)

func New() pusher.Stat {
	return &Ding{}
}

func (d *Ding) Push(msgType, text string, extra ...string) error {
	msg := msg{MsgType: msgType}
	if len(extra) > 0 {
		extraMsg := strings.Join(extra, "\r\n")
		text = text + "\r\n" + extraMsg
	}
	msg.Text.Content = text
	msgBytes, _ := json.Marshal(msg)
	res, err := utiles.PostWithJson(msgBytes, AccessUrl)
	if err != nil {
		return err
	}
	var r map[string]interface{}
	if _ = json.Unmarshal(res, &r); r != nil {
		if r["errcode"].(float64) == 0 {
			return nil
		}
	}
	return errors.New("op fail")
}
