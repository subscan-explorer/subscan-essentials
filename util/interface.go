package util

import (
	"encoding/json"
	"runtime"
	"strings"
)

// Func
func CallerName() string {
	pc, _, _, _ := runtime.Caller(2)
	return runtime.FuncForPC(pc).Name()
}

func GetFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}

// Interface
func BoolFromInterface(i interface{}) bool {
	switch i := i.(type) {
	case string:
		return strings.ToLower(i) == "true"
	case bool:
		return i
	}
	return false
}

func ToString(i interface{}) string {
	var val string
	switch i := i.(type) {
	case string:
		val = i
	case []byte:
		val = string(i)
	default:
		b, _ := json.Marshal(i)
		val = string(b)
	}
	return val
}

func UnmarshalAny(r interface{}, raw interface{}) {
	switch raw := raw.(type) {
	case string:
		_ = json.Unmarshal([]byte(raw), &r)
	case []uint8:
		_ = json.Unmarshal(raw, &r)
	default:
		b, _ := json.Marshal(raw)
		_ = json.Unmarshal(b, r)
	}
}
