package tests

import (
	"github.com/itering/subscan/util"
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
	funcNameRaw := "github.com/itering/subscan/tests.TestFunc"

	callerName := util.CallerName()
	funcName := util.GetFuncName()

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

	if util.BoolFromInterface(boolean) == !boolean {
		t.Errorf(
			"Parse bool failed, got %v, want %v",
			!boolean,
			boolean,
		)
	}

	if util.BoolFromInterface(rt) == false {
		t.Errorf(
			"Parse bool failed, got %v, want %v",
			rf,
			rt,
		)
	}

	if util.BoolFromInterface(rf) == true {
		t.Errorf(
			"Parse bool failed, got %v, want %v",
			rt,
			rf,
		)
	}
}

func TestFieldName(t *testing.T) {
	a, _ := util.GetStringValueByFieldName(field, "A")
	b, _ := util.GetStringValueByFieldName(field, "B")

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
