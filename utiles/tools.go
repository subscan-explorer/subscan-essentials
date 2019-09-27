package utiles

import (
	"encoding/hex"
	"fmt"
	"github.com/shopspring/decimal"
	"math/big"
	"strconv"
	"strings"
	"time"
)

func IntToString(i int) string {
	return strconv.Itoa(i)
}

func StringToInt(s string) int {
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	return 0
}

func InsertInts(o []int, index int, new int) []int {
	if index > len(o) {
		return append(o, new)
	}
	temp := append([]int{new}, o[index:]...)
	return append(o[:index], temp...)
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func IntInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
func AddHex(s string) string {
	if strings.HasPrefix(s, "0x") {
		return s
	}
	return strings.ToLower("0x" + s)
}

func U256(v string) *big.Int {
	v = strings.TrimPrefix(v, "0x")
	bn := new(big.Int)
	n, _ := bn.SetString(v, 16)
	return n
}

func HexToNumStr(v string) string {
	return U256(v).String()
}

func HexToBytes(s string) []byte {
	s = strings.TrimPrefix(s, "0x")
	c := make([]byte, hex.DecodedLen(len(s)))
	_, _ = hex.Decode(c, []byte(s))
	return c
}

func BytesToHex(b []byte) string {
	c := make([]byte, hex.EncodedLen(len(b)))
	hex.Encode(c, b)
	return string(c)
}

func BigToDecimal(v *big.Int) decimal.Decimal {
	return decimal.NewFromBigInt(v, 0).Div(decimal.NewFromFloat(1E9))
}

func FloatToDecimal(v float64) decimal.Decimal {
	bigint, _ := big.NewFloat(v).Int(big.NewInt(0))
	return BigToDecimal(bigint)
}

func IntToHex(i interface{}) string {
	return fmt.Sprintf("%x", i)
}

func TrimHex(s string) string {
	return strings.TrimPrefix(s, "0x")
}

func ContinuousSlice(start, count int, order string) (r []int) {
	if count <= 0 || start <= 0 {
		return
	}
	for i := 0; i < count; i++ {
		if order == "desc" {
			if start-i < 0 {
				break
			}
			r = append(r, start-i)
		} else {
			r = append(r, start+i)
		}
	}
	return
}

func UTCtoTimestamp(utcTime string) int {
	if strings.HasSuffix(utcTime, "Z") == false {
		utcTime = utcTime + "Z"
	}
	if t, err := time.Parse(time.RFC3339, utcTime); err == nil {
		return int(t.Unix())
	}
	return 0
}
