package model

import "encoding/json"

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
