package util

import (
	"testing"
)

var (
	field = struct {
		A string
		B string
	}{
		"a",
		"b",
	}
)

func TestFunc(t *testing.T) {
	callerNameRaw := "testing.tRunner"
	funcNameRaw := "github.com/itering/subscan/internal/util.TestFunc"

	callerName := CallerName()
	funcName := GetFuncName()

	if funcName != funcNameRaw {
		t.Errorf(
			"Get function name failed, got %s, want %s",
			funcName,
			funcNameRaw,
		)
	}

	if callerName != callerNameRaw {
		t.Errorf(
			"Get caller name failed, got %s, want %s",
			callerName,
			callerNameRaw,
		)
	}
}

func TestBool(t *testing.T) {
	boolean := false
	rt := "true"
	rf := "false"

	if BoolFromInterface(boolean) == !boolean {
		t.Errorf(
			"Parse bool failed, got %v, want %v",
			!boolean,
			boolean,
		)
	}

	if BoolFromInterface(rt) == false {
		t.Errorf(
			"Parse bool failed, got %v, want %v",
			rf,
			rt,
		)
	}

	if BoolFromInterface(rf) == true {
		t.Errorf(
			"Parse bool failed, got %v, want %v",
			rt,
			rf,
		)
	}
}

func TestFieldName(t *testing.T) {
	a, _ := GetStringValueByFieldName(field, "A")
	b, _ := GetStringValueByFieldName(field, "B")

	if a != "a" {
		t.Errorf(
			"Get struct field a failed, got %v, want %v",
			a,
			"a",
		)
	}

	if b != "b" {
		t.Errorf(
			"Get struct field a failed, got %v, want %v",
			b,
			"b",
		)
	}
}
