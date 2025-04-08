package util

import (
	"cmp"
	"encoding/binary"
	"math/big"
	"sort"
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
	if v == "" {
		return big.NewInt(0)
	}
	bn := new(big.Int)
	n, _ := bn.SetString(v, 16)
	return n
}

func DecimalFromU256(v string) decimal.Decimal {
	if u256 := U256(v); u256 != nil {
		return decimal.NewFromBigInt(U256(v), 0)
	}
	return decimal.Zero
}

// IntFromInterface Convert int64, uint64, float64, string to int, return 0 if other types
func IntFromInterface(i any) int {
	switch i := i.(type) {
	case int8:
		return int(i)
	case int16:
		return int(i)
	case int32:
		return int(i)
	case int64:
		return int(i)
	case int:
		return i
	case uint8:
		return int(i)
	case uint16:
		return int(i)
	case uint32:
		return int(i)
	case uint64:
		return int(i)
	case uint:
		return int(i)
	case float64:
		return int(i)
	case float32:
		return int(i)
	case string:
		return StringToInt(i)
	}
	return 0
}

func UIntFromInterface(i interface{}) uint {
	switch i := i.(type) {
	case int8:
		return uint(i)
	case int16:
		return uint(i)
	case int32:
		return uint(i)
	case int64:
		return uint(i)
	case int:
		return uint(i)
	case uint8:
		return uint(i)
	case uint16:
		return uint(i)
	case uint32:
		return uint(i)
	case uint64:
		return uint(i)
	case uint:
		return i
	case float64:
		return uint(i)
	case float32:
		return uint(i)
	case string:
		if i, err := strconv.ParseUint(i, 10, 0); err == nil {
			return uint(i)
		}
		return 0
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
	case uint32:
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
	case uint:
		return decimal.NewFromUint64(uint64(i))
	case string:
		r, _ := decimal.NewFromString(i)
		return r
	case decimal.Decimal:
		return i
	case *big.Int:
		return decimal.NewFromBigInt(i, 0)
	}
	return decimal.Zero
}

// Big Int
func U32Encode(i uint32) string {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, i)
	return BytesToHex(bs)
}

func U32Decode(v string) uint32 {
	return binary.LittleEndian.Uint32(HexToBytes(v))
}

func U16Encode(i uint16) string {
	bs := make([]byte, 2)
	binary.LittleEndian.PutUint16(bs, i)
	return BytesToHex(bs)
}

func U8Encode(i int) string {
	bs := make([]byte, 1)
	bs[0] = byte(i)
	return BytesToHex(bs)
}

func BigIntFromInterface(i interface{}) *big.Int {
	switch i := i.(type) {
	case int:
		return big.NewInt(int64(i))
	case int64:
		return big.NewInt(i)
	case float64:
		return big.NewInt(int64(i))
	case string:
		b := big.NewInt(0)
		b.SetString(i, 10)
		return b
	}
	return big.NewInt(0)
}

// evm u256 type convert
func EvmReverseU256Decoder(u256 string) decimal.Decimal {
	reverseData := Reverse(HexToBytes(u256))
	return decimal.NewFromBigInt(U256(BytesToHex(reverseData.([]byte))), 0)
}

func EvmReverseU256DecoderToBigInt(u256 string) *big.Int {
	reverseData := Reverse(HexToBytes(u256))
	return U256(BytesToHex(reverseData.([]byte)))
}

func EvmU256Decoder(u256 string) decimal.Decimal {
	return decimal.NewFromBigInt(U256(u256), 0)
}

func Min[T cmp.Ordered](x, y T) T {
	if x < y {
		return x
	}
	return y
}

func FillMissedInt(start, end int, exist []int) []int {
	sort.Slice(exist, func(i, j int) bool { return exist[i] < exist[j] })
	var missedBlocks []int
	var idx = 0
	cur := start
	for idx < len(exist) {
		if exist[idx] > end {
			break
		}
		if cur == exist[idx] {
			cur++
			idx++
		} else if cur < exist[idx] {
			missedBlocks = append(missedBlocks, cur)
			cur++
		} else {
			idx++
		}
	}
	if len(exist) != 0 {
		for i := exist[len(exist)-1] + 1; i <= end; i++ {
			missedBlocks = append(missedBlocks, i)
		}
	} else {
		for i := start; i <= end; i++ {
			missedBlocks = append(missedBlocks, i)
		}
	}
	return missedBlocks
}

func StringToUInt(s string) uint {
	if i, err := strconv.ParseUint(s, 10, 64); err == nil {
		return uint(i)
	}
	return 0
}
