package model

import (
	"encoding/json"
	"time"

	"github.com/itering/subscan/util/es"
)

var esClient *es.EsClient

func InitEsClient() {
	if esClient == nil {
		esClient, _ = es.NewEsClient()
	}
}

type Subscriber struct {
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	Email     string    `gorm:"type:varchar(64);" json:"email" `
}

func ParsingEventParam(params interface{}) (param []EventParam, err error) {
	var data []byte
	switch p := params.(type) {
	case []uint8:
		data = p
	case []interface{}:
		data, err = json.Marshal(p)
		if err != nil {
			return
		}
	case string:
		data = []byte(p)
	}
	err = json.Unmarshal(data, &param)
	if err != nil {
		return
	}
	return
}

func ParsingExtrinsicParam(params interface{}) (param []ExtrinsicParam) {
	var data []byte
	var err error
	switch p := params.(type) {
	case []uint8:
		data = p
	case string:
		data = []byte(p)
	default:
		data, err = json.Marshal(p)
		if err != nil {
			return
		}
	}
	err = json.Unmarshal(data, &param)
	if err != nil {
		return
	}
	return
}

func ParsingExtrinsicErrorParam(params interface{}) (param map[string]interface{}) {
	var data []byte
	var err error
	switch p := params.(type) {
	case []uint8:
		data = p
	case string:
		data = []byte(p)
	default:
		data, err = json.Marshal(p)
		if err != nil {
			return
		}
	}
	err = json.Unmarshal(data, &param)
	if err != nil {
		return
	}
	return
}

type Call struct {
	CallFunction string           `json:"call_function"`
	CallModule   string           `json:"call_module"`
	CallIndex    string           `json:"call_index"`
	CallArgs     []ExtrinsicParam `json:"call_args"`
}

type BoxProposal struct {
	CallModule string           `json:"call_module"`
	CallName   string           `json:"call_name"`
	Params     []ExtrinsicParam `json:"params"`
	CallIndex  string           `json:"call_index"`
}
