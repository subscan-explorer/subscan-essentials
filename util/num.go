package util

import (
	"math/big"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

// Int
func IntToString(i int) string {
	return strconv.Itoa(i)
}

// Convert str to int, return 0 if error
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

func IntInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func U256(v string) *big.Int {
	v = strings.TrimPrefix(v, "0x")
	bn := new(big.Int)
	n, _ := bn.SetString(v, 16)
	return n
}

// Convert int64, uint64, float64, string to int, return 0 if other types
func IntFromInterface(i interface{}) int {
	switch i := i.(type) {
	case int:
		return i
	case int64:
		return int(i)
	case uint64:
		return int(i)
	case float64:
		return int(i)
	case string:
		return StringToInt(i)
	}
	return 0
}

func Int64FromInterface(i interface{}) int64 {
	switch i := i.(type) {
	case int:
		return int64(i)
	case int64:
		return i
	case uint64:
		return int64(i)
	case float64:
		return int64(i)
	case string:
		r, _ := strconv.ParseInt(i, 10, 64)
		return r
	}
	return 0
}

func DecimalFromInterface(i interface{}) decimal.Decimal {
	switch i := i.(type) {
	case int:
		return decimal.New(int64(i), 0)
	case int64:
		return decimal.New(i, 0)
	case uint64:
		return decimal.New(int64(i), 0)
	case float64:
		return decimal.NewFromFloat(i)
	case string:
		r, _ := decimal.NewFromString(i)
		return r
	}
	return decimal.Zero
}
