package util

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"runtime"
	"strconv"
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

func UnmarshalAny(r interface{}, raw interface{}) (err error) {
	switch raw := raw.(type) {
	case string:
		err = json.Unmarshal([]byte(raw), &r)
	case []uint8:
		err = json.Unmarshal(raw, &r)
	default:
		b, _ := json.Marshal(raw)
		err = json.Unmarshal(b, r)
	}
	return err
}

// ToBytes convert interface to bytes
// support string, []byte, struct, map, slice.
// PLEASE DO NOT USE IT FOR int, float, bool
func ToBytes(i interface{}) (bytesData []byte) {
	switch v := i.(type) {
	case string:
		bytesData = []byte(v)
	case []byte:
		bytesData = v
	default:
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		bytesData, _ = json.Marshal(i)
	}
	return bytesData
}

func ToString(i interface{}) string {
	var val string
	switch i := i.(type) {
	case string:
		val = i
	case int64:
		val = strconv.FormatInt(i, 10)
	case int:
		val = strconv.FormatInt(int64(i), 10)
	case uint:
		val = strconv.FormatUint(uint64(i), 10)
	case uint64:
		val = strconv.FormatUint(i, 10)
	case []byte:
		val = string(i)
	case []int:
		var bytes []byte
		for _, b := range i {
			bytes = append(bytes, byte(b))
		}
		val = string(bytes)
	case float64:
		return fmt.Sprintf("%f", i)
	default:
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		b, _ := json.Marshal(i)
		val = string(b)
		if val == "null" {
			return ""
		}
	}
	return val
}

func Base64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func Base64Decode(s string) string {
	d, _ := base64.StdEncoding.DecodeString(s)
	return string(d)
}
