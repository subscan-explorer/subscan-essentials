package tests

import (
	"github.com/itering/subscan/util"

	"strconv"
	"testing"
)

var (
	prefixTests = []struct {
		raw    string
		prefix string
		trim   string
	}{
		{"", "", ""},
		{" ", "", ""},
		{"0x", "0x", ""},
		{"00", "0x00", "00"},
	}
)

func TestPrefix(t *testing.T) {
	for x, test := range prefixTests {
		res := util.AddHex(test.raw)
		if res != test.prefix {
			t.Errorf(
				"Add Hex #%d failed, got: %s want: %s",
				x, res, test.prefix,
			)
		}

		if trim := util.TrimHex(res); trim != test.trim {
			t.Errorf(
				"Trim Hex #%d failed, got: %s want: %s",
				x, res, test.trim,
			)
		}
	}
}

func TestNums(t *testing.T) {
	for i := 16; i < 99; i++ {
		str := strconv.Itoa(i)
		hex := util.IntToHex(i)
		bytes := util.HexToBytes(hex)
		bHex := util.BytesToHex(bytes)
		hStr := util.HexToNumStr(hex)
		if hex != bHex {
			t.Errorf("%x", bytes)
			t.Errorf(
				"Convert Hex #%d failed, got: %s want: %s",
				i, bHex, hex,
			)
		}

		if str != hStr {
			t.Errorf(
				"Convert Hex #%d failed, got: %s want: %s",
				i, hStr, str,
			)
		}
	}
}
