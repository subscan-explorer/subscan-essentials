package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFunc(t *testing.T) {
	callerNameRaw := "testing.tRunner"
	funcNameRaw := "github.com/itering/subscan/util.TestFunc"

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
	rint := 1

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

	if BoolFromInterface(rint) == true {
		t.Errorf(
			"Parse bool failed, got %v, want %v",
			rt,
			rf,
		)
	}
}

func TestToString(t *testing.T) {
	testCase := []struct {
		i interface{}
		r string
	}{
		{"abc", "abc"},
		{[]byte{97, 98, 99}, "abc"},
		{map[string]string{"a": "b"}, `{"a":"b"}`},
	}
	for _, test := range testCase {
		assert.Equal(t, test.r, ToString(test.i))
	}
}

func TestUnmarshalAny(t *testing.T) {
	p := new(struct {
		One int
		Two int
	})
	mapTest := map[string]int{"one": 1, "two": 2}
	UnmarshalAny(&p, mapTest)
	assert.Equal(t, &struct {
		One int
		Two int
	}{1, 2}, p)

	UnmarshalAny(&p, `{"one":21,"two":22}`)
	assert.Equal(t, &struct {
		One int
		Two int
	}{21, 22}, p)

	UnmarshalAny(&p, []byte(`{"one":31,"two":32}`))
	assert.Equal(t, &struct {
		One int
		Two int
	}{31, 32}, p)

}
