package dingding

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/itering/subscan/util"
	"github.com/itering/subscan/util/pusher"
)

const (
	AccessUrl      = "https://oapi.dingtalk.com/robot/send?access_token=28b1fcc0d0e5da3c8e9178fbf3aa011dbbf750f823a584729869262a206acc27"
	AlertAccessUrl = "https://oapi.dingtalk.com/robot/send?access_token=340eb961c25046fda28476345205dda7e0e21a21a40a120eaf8db32379c593f8"
)

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

func New() pusher.Alert {
	return &Ding{}
}

func (d *Ding) Push(scene, text string, extra ...string) error {
	msg := msg{
		MsgType: "text",
	}
	msg.Text.Content = fmt.Sprintf("%s %s", util.HostName, scene) + "\r\n" + text + "\r\n" + strings.Join(extra, "\r\n")
	msgBytes, _ := json.Marshal(msg)
	if scene == "Subscan" {
		_, _ = util.PostWithJson(msgBytes, AccessUrl)
		return nil
	}
	_, _ = util.PostWithJson(msgBytes, AlertAccessUrl)
	return nil
}
