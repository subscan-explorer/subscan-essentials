package util

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"

	"golang.org/x/exp/slog"
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

func StringFromInterface(i interface{}) (string, error) {
	switch i := i.(type) {
	case string:
		return i, nil
	case []byte:
		return string(i), nil
	}
	return "", fmt.Errorf("error converting interface to string. value: %+v, type: %T", i, i)
}

func FromString(s string, out *interface{}) {
	switch (*out).(type) {
	case string:
		*out = s
	case []byte:
		*out = []byte(s)
	default:
		bytes := []byte(s)
		err := json.Unmarshal(bytes, out)
		if err != nil {
			slog.Warn("error unmarshalling string to interface", "string", s, "error", err)
		}
	}
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
