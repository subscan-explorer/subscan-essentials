package util

import (
	"encoding/json"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
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

func InterfaceToString(i interface{}) string {
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

func GetStringValueByFieldName(n interface{}, fieldName string) (string, bool) {
	s := reflect.ValueOf(n)
	if s.Kind() == reflect.Ptr {
		s = s.Elem()
	}
	if s.Kind() != reflect.Struct {
		return "", false
	}
	f := s.FieldByName(fieldName)
	if !f.IsValid() {
		return "", false
	}
	switch f.Kind() {
	case reflect.String:
		return f.Interface().(string), true
	case reflect.Int:
		return strconv.FormatInt(f.Int(), 10), true
	case reflect.Uint:
		return strconv.FormatUint(f.Uint(), 10), true
	case reflect.Int64:
		return strconv.FormatInt(f.Int(), 10), true
	default:
		if reflect.TypeOf(f.Interface()).Name() == "Decimal" {
			return f.Interface().(decimal.Decimal).String(), true
		}
		return "", false
	}
}

func UnmarshalToAnything(r interface{}, raw interface{}) {
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
