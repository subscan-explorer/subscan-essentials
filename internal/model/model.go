package model

import "encoding/json"

func ParsingEventParam(params interface{}) (param []EventParam, err error) {
	var data []byte
	switch params.(type) {
	case []uint8:
		data = params.([]uint8)
	case []interface{}:
		data, err = json.Marshal(params)
		if err != nil {
			return
		}
	case string:
		data = []byte(params.(string))
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
	switch params.(type) {
	case []uint8:
		data = params.([]uint8)
	case string:
		data = []byte(params.(string))
	default:
		data, err = json.Marshal(params)
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
	switch params.(type) {
	case []uint8:
		data = params.([]uint8)
	case string:
		data = []byte(params.(string))
	default:
		data, err = json.Marshal(params)
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
