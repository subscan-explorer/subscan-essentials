package dingding

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/itering/subscan/util"
	"github.com/itering/subscan/util/pusher"
)

var DingClient = New()

type (
	Ding struct {
		AccessUrl      string
		AlertAccessUrl string
	}
	msg struct {
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
		_, _ = util.PostWithJson(msgBytes, d.AccessUrl)
		return nil
	}
	_, _ = util.PostWithJson(msgBytes, d.AlertAccessUrl)
	return nil
}
