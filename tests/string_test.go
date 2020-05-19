package tests

import (
	"testing"

	"github.com/itering/subscan/util"
)

var (
	notCamel   = "camelcase"
	camel      = "Camelcase"
	upperCamel = "Camelcase"
	alice      = []string{"a", "b", "c", "d"}
	bob        = []string{"c", "d", "e", "f"}
	abi        = []string{"c", "d"}
	abe        = []string{"a", "b"}
)

func TestCamel(t *testing.T) {
	// TODO
	//
	// This is not as expected
	ucamel := util.CamelString(notCamel)
	if ucamel != camel {
		t.Errorf(
			"Camel string failed, got: %s, want: %s",
			ucamel,
			camel,
		)
	}

	uUpperCamel := util.UpperCamel(ucamel)
	if uUpperCamel != upperCamel {
		t.Errorf(
			"Camel string failed, got: %s, want: %s",
			uUpperCamel,
			upperCamel,
		)
	}
}

func TestSets(t *testing.T) {
	uAbi := util.StringsIntersection(alice, bob)
	abiLen := len(abi)
	uAbiLen := len(uAbi)
	if abiLen != uAbiLen {
		t.Errorf(
			"Map string to string slice length failed, got %v, want %v",
			abiLen,
			uAbiLen,
		)
	}

	for i := range uAbi {
		if !util.StringInSlice(uAbi[i], abi) {
			t.Errorf(
				"Get string intersection failed #%d, got %v, want %v",
				i,
				uAbi[i],
				abi[i],
			)
		}
	}

	uAbe := util.StringsExclude(alice, bob)
	for i := range uAbe {
		if !util.StringInSlice(uAbe[i], abe) {
			t.Errorf(
				"Get string exclude failed #%d, got %v, want %v",
				i,
				uAbe[i],
				abe[i],
			)
		}
	}
}
