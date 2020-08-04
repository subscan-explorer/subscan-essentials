package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
	emptyStr := ""
	assert.Equal(t, "", CamelString(emptyStr))
	assert.Equal(t, "", UpperCamel(emptyStr))

	ucamel := CamelString(notCamel)
	if ucamel != camel {
		t.Errorf(
			"Camel string failed, got: %s, want: %s",
			ucamel,
			camel,
		)
	}

	uUpperCamel := UpperCamel(ucamel)
	if uUpperCamel != upperCamel {
		t.Errorf(
			"Camel string failed, got: %s, want: %s",
			uUpperCamel,
			upperCamel,
		)
	}
}

func TestSets(t *testing.T) {
	uAbi := StringsIntersection(alice, bob)
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
		if !StringInSlice(uAbi[i], abi) {
			t.Errorf(
				"Get string intersection failed #%d, got %v, want %v",
				i,
				uAbi[i],
				abi[i],
			)
		}
	}

	uAbe := StringsExclude(alice, bob)
	for i := range uAbe {
		if !StringInSlice(uAbe[i], abe) {
			t.Errorf(
				"Get string exclude failed #%d, got %v, want %v",
				i,
				uAbe[i],
				abe[i],
			)
		}
	}
}
