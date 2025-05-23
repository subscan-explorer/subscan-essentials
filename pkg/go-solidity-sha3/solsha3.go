package solsha3

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/itering/scale.go/pkg/go-ethereum/crypto/sha3"
	"math/big"
	"reflect"
)

// Lengths of hashes and addresses in bytes.
const (
	// HashLength is the expected length of the hash
	HashLength = 32
	// AddressLength is the expected length of the address
	AddressLength = 20
)

// Address represents the 20 byte address of an Ethereum account.
type Address [AddressLength]byte

func BytesToAddress(b []byte) Address {
	var a Address
	a.SetBytes(b)
	return a
}

func (a *Address) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressLength:]
	}
	copy(a[AddressLength-len(b):], b)
}

// Bytes gets the string representation of the underlying address.
func (a Address) Bytes() []byte { return a[:] }

// Address address
func ToAddress(input interface{}) []byte {
	switch v := input.(type) {
	case string:
		return HexToAddress(v).Bytes()[:]
	default:
		return HexToAddress("").Bytes()[:]
	}
}

// Uint256 uint256
func Uint256(input interface{}) []byte {
	switch v := input.(type) {
	case *big.Int:
		return U256(v)
	case string:
		bn := new(big.Int)
		bn.SetString(v, 10)
		return U256(bn)
	default:
		return RightPadBytes([]byte(""), 32)
	}
}

// Uint128 uint128
func Uint128(input interface{}) []byte {
	switch v := input.(type) {
	case *big.Int:
		return LeftPadBytes(v.Bytes(), 16)
	case string:
		bn := new(big.Int)
		bn.SetString(v, 10)
		return LeftPadBytes(bn.Bytes(), 16)
	default:
		return LeftPadBytes([]byte(""), 16)
	}
}

// Uint64 uint64
func Uint64(input interface{}) []byte {
	b := new(bytes.Buffer)
	switch v := input.(type) {
	case *big.Int:
		binary.Write(b, binary.BigEndian, v.Uint64())
	case string:
		bn := new(big.Int)
		bn.SetString(v, 10)
		binary.Write(b, binary.BigEndian, bn.Uint64())
	case uint64:
		binary.Write(b, binary.BigEndian, v)
	case uint32:
		binary.Write(b, binary.BigEndian, uint64(v))
	case uint16:
		binary.Write(b, binary.BigEndian, uint64(v))
	case uint8:
		binary.Write(b, binary.BigEndian, uint64(v))
	case uint:
		binary.Write(b, binary.BigEndian, uint64(v))
	case int64:
		binary.Write(b, binary.BigEndian, uint64(v))
	case int32:
		binary.Write(b, binary.BigEndian, uint64(v))
	case int16:
		binary.Write(b, binary.BigEndian, uint64(v))
	case int8:
		binary.Write(b, binary.BigEndian, uint64(v))
	case int:
		binary.Write(b, binary.BigEndian, uint64(v))
	default:
		binary.Write(b, binary.BigEndian, uint64(0))
	}
	return b.Bytes()
}

// Uint32 uint32
func Uint32(input interface{}) []byte {
	b := new(bytes.Buffer)
	switch v := input.(type) {
	case *big.Int:
		binary.Write(b, binary.BigEndian, uint32(v.Uint64()))
	case string:
		bn := new(big.Int)
		bn.SetString(v, 10)
		binary.Write(b, binary.BigEndian, uint32(bn.Uint64()))
	case uint64:
		binary.Write(b, binary.BigEndian, uint32(v))
	case uint32:
		binary.Write(b, binary.BigEndian, uint32(v))
	case uint16:
		binary.Write(b, binary.BigEndian, uint32(v))
	case uint8:
		binary.Write(b, binary.BigEndian, uint32(v))
	case uint:
		binary.Write(b, binary.BigEndian, uint32(v))
	case int64:
		binary.Write(b, binary.BigEndian, uint32(v))
	case int32:
		binary.Write(b, binary.BigEndian, v)
	case int16:
		binary.Write(b, binary.BigEndian, uint32(v))
	case int8:
		binary.Write(b, binary.BigEndian, uint32(v))
	case int:
		binary.Write(b, binary.BigEndian, uint32(v))
	default:
		binary.Write(b, binary.BigEndian, uint32(0))
	}
	return b.Bytes()
}

// Uint16 uint16
func Uint16(input interface{}) []byte {
	b := new(bytes.Buffer)
	switch v := input.(type) {
	case *big.Int:
		binary.Write(b, binary.BigEndian, uint16(v.Uint64()))
	case string:
		bn := new(big.Int)
		bn.SetString(v, 10)
		binary.Write(b, binary.BigEndian, uint16(bn.Uint64()))
	case uint64:
		binary.Write(b, binary.BigEndian, uint16(v))
	case uint32:
		binary.Write(b, binary.BigEndian, uint16(v))
	case uint16:
		binary.Write(b, binary.BigEndian, v)
	case uint8:
		binary.Write(b, binary.BigEndian, uint16(v))
	case uint:
		binary.Write(b, binary.BigEndian, uint16(v))
	case int64:
		binary.Write(b, binary.BigEndian, uint16(v))
	case int32:
		binary.Write(b, binary.BigEndian, uint16(v))
	case int16:
		binary.Write(b, binary.BigEndian, uint16(v))
	case int8:
		binary.Write(b, binary.BigEndian, uint16(v))
	case int:
		binary.Write(b, binary.BigEndian, uint16(v))
	default:
		binary.Write(b, binary.BigEndian, uint16(0))
	}
	return b.Bytes()
}

// Uint8 uint8
func Uint8(input interface{}) []byte {
	b := new(bytes.Buffer)
	switch v := input.(type) {
	case *big.Int:
		binary.Write(b, binary.BigEndian, uint8(v.Uint64()))
	case string:
		bn := new(big.Int)
		bn.SetString(v, 10)
		binary.Write(b, binary.BigEndian, uint8(bn.Uint64()))
	case uint64:
		binary.Write(b, binary.BigEndian, uint8(v))
	case uint32:
		binary.Write(b, binary.BigEndian, uint8(v))
	case uint16:
		binary.Write(b, binary.BigEndian, uint8(v))
	case uint8:
		binary.Write(b, binary.BigEndian, v)
	case uint:
		binary.Write(b, binary.BigEndian, uint8(v))
	case int64:
		binary.Write(b, binary.BigEndian, uint8(v))
	case int32:
		binary.Write(b, binary.BigEndian, uint8(v))
	case int16:
		binary.Write(b, binary.BigEndian, uint8(v))
	case int8:
		binary.Write(b, binary.BigEndian, uint8(v))
	case int:
		binary.Write(b, binary.BigEndian, uint8(v))
	default:
		binary.Write(b, binary.BigEndian, uint8(0))
	}
	return b.Bytes()
}

// Int256 int256
func Int256(input interface{}) []byte {
	switch v := input.(type) {
	case *big.Int:
		return LeftPadBytes(v.Bytes(), 32)
	case string:
		bn := new(big.Int)
		bn.SetString(v, 10)
		return LeftPadBytes(bn.Bytes(), 32)
	case uint64:
		bn := big.NewInt(int64(v))
		return LeftPadBytes(bn.Bytes(), 32)
	case uint32:
		bn := big.NewInt(int64(v))
		return LeftPadBytes(bn.Bytes(), 32)
	case uint16:
		bn := big.NewInt(int64(v))
		return LeftPadBytes(bn.Bytes(), 32)
	case uint8:
		bn := big.NewInt(int64(v))
		return LeftPadBytes(bn.Bytes(), 32)
	case uint:
		bn := big.NewInt(int64(v))
		return LeftPadBytes(bn.Bytes(), 32)
	case int64:
		bn := big.NewInt(int64(v))
		return LeftPadBytes(bn.Bytes(), 32)
	case int32:
		bn := big.NewInt(int64(v))
		return LeftPadBytes(bn.Bytes(), 32)
	case int16:
		bn := big.NewInt(int64(v))
		return LeftPadBytes(bn.Bytes(), 32)
	case int8:
		bn := big.NewInt(int64(v))
		return LeftPadBytes(bn.Bytes(), 32)
	case int:
		bn := big.NewInt(int64(v))
		return LeftPadBytes(bn.Bytes(), 32)
	default:
		bn := big.NewInt(int64(0))
		return LeftPadBytes(bn.Bytes(), 32)
	}
}

// Int128 int128
func Int128(input interface{}) []byte {
	switch v := input.(type) {
	case *big.Int:
		return LeftPadBytes(v.Bytes(), 16)
	case string:
		bn := new(big.Int)
		bn.SetString(v, 10)
		return LeftPadBytes(bn.Bytes(), 16)
	case uint64:
		bn := big.NewInt(int64(v))
		return LeftPadBytes(bn.Bytes(), 16)
	case uint32:
		bn := big.NewInt(int64(v))
		return LeftPadBytes(bn.Bytes(), 16)
	case uint16:
		bn := big.NewInt(int64(v))
		return LeftPadBytes(bn.Bytes(), 16)
	case uint8:
		bn := big.NewInt(int64(v))
		return LeftPadBytes(bn.Bytes(), 16)
	case uint:
		bn := big.NewInt(int64(v))
		return LeftPadBytes(bn.Bytes(), 16)
	case int64:
		bn := big.NewInt(int64(v))
		return LeftPadBytes(bn.Bytes(), 16)
	case int32:
		bn := big.NewInt(int64(v))
		return LeftPadBytes(bn.Bytes(), 16)
	case int16:
		bn := big.NewInt(int64(v))
		return LeftPadBytes(bn.Bytes(), 16)
	case int8:
		bn := big.NewInt(int64(v))
		return LeftPadBytes(bn.Bytes(), 16)
	case int:
		bn := big.NewInt(int64(v))
		return LeftPadBytes(bn.Bytes(), 16)
	default:
		bn := big.NewInt(int64(0))
		return LeftPadBytes(bn.Bytes(), 16)
	}
}

// Int64 int64
func Int64(input interface{}) []byte {
	b := make([]byte, 8)
	switch v := input.(type) {
	case *big.Int:
		binary.BigEndian.PutUint64(b, v.Uint64())
	case string:
		bn := new(big.Int)
		bn.SetString(v, 10)
		binary.BigEndian.PutUint64(b, bn.Uint64())
	case uint64:
		binary.BigEndian.PutUint64(b, v)
	case uint32:
		binary.BigEndian.PutUint64(b, uint64(v))
	case uint16:
		binary.BigEndian.PutUint64(b, uint64(v))
	case uint8:
		binary.BigEndian.PutUint64(b, uint64(v))
	case uint:
		binary.BigEndian.PutUint64(b, uint64(v))
	case int64:
		binary.BigEndian.PutUint64(b, uint64(v))
	case int32:
		binary.BigEndian.PutUint64(b, uint64(v))
	case int16:
		binary.BigEndian.PutUint64(b, uint64(v))
	case int8:
		binary.BigEndian.PutUint64(b, uint64(v))
	case int:
		binary.BigEndian.PutUint64(b, uint64(v))
	default:
		binary.BigEndian.PutUint64(b, uint64(0))
	}
	return b
}

// Int32 int32
func Int32(input interface{}) []byte {
	b := make([]byte, 4)
	switch v := input.(type) {
	case *big.Int:
		binary.BigEndian.PutUint32(b, uint32(v.Uint64()))
	case string:
		bn := new(big.Int)
		bn.SetString(v, 10)
		binary.BigEndian.PutUint32(b, uint32(bn.Uint64()))
	case uint64:
		binary.BigEndian.PutUint32(b, uint32(v))
	case uint32:
		binary.BigEndian.PutUint32(b, v)
	case uint16:
		binary.BigEndian.PutUint32(b, uint32(v))
	case uint8:
		binary.BigEndian.PutUint32(b, uint32(v))
	case uint:
		binary.BigEndian.PutUint32(b, uint32(v))
	case int64:
		binary.BigEndian.PutUint32(b, uint32(v))
	case int32:
		binary.BigEndian.PutUint32(b, uint32(v))
	case int16:
		binary.BigEndian.PutUint32(b, uint32(v))
	case int8:
		binary.BigEndian.PutUint32(b, uint32(v))
	case int:
		binary.BigEndian.PutUint32(b, uint32(v))
	default:
		binary.BigEndian.PutUint32(b, uint32(0))
	}
	return b
}

// Int16 int16
func Int16(input interface{}) []byte {
	b := make([]byte, 2)
	switch v := input.(type) {
	case *big.Int:
		binary.BigEndian.PutUint16(b, uint16(v.Uint64()))
	case string:
		bn := new(big.Int)
		bn.SetString(v, 10)
		binary.BigEndian.PutUint16(b, uint16(bn.Uint64()))
	case uint64:
		binary.BigEndian.PutUint16(b, uint16(v))
	case uint32:
		binary.BigEndian.PutUint16(b, uint16(v))
	case uint16:
		binary.BigEndian.PutUint16(b, v)
	case uint8:
		binary.BigEndian.PutUint16(b, uint16(v))
	case uint:
		binary.BigEndian.PutUint16(b, uint16(v))
	case int64:
		binary.BigEndian.PutUint16(b, uint16(v))
	case int32:
		binary.BigEndian.PutUint16(b, uint16(v))
	case int16:
		binary.BigEndian.PutUint16(b, uint16(v))
	case int8:
		binary.BigEndian.PutUint16(b, uint16(v))
	case int:
		binary.BigEndian.PutUint16(b, uint16(v))
	default:
		binary.BigEndian.PutUint16(b, uint16(0))
	}
	return b
}

// Int8 int8
func Int8(input interface{}) []byte {
	b := make([]byte, 1)
	switch v := input.(type) {
	case *big.Int:
		b[0] = byte(int8(v.Uint64()))
	case string:
		bn := new(big.Int)
		bn.SetString(v, 10)
		b[0] = byte(int8(bn.Uint64()))
	case uint64:
		b[0] = byte(int8(v))
	case uint32:
		b[0] = byte(int8(v))
	case uint16:
		b[0] = byte(int8(v))
	case uint8:
		b[0] = byte(int8(v))
	case uint:
		b[0] = byte(int8(v))
	case int64:
		b[0] = byte(int8(v))
	case int32:
		b[0] = byte(int8(v))
	case int16:
		b[0] = byte(int8(v))
	case int8:
		b[0] = byte(v)
	case int:
		b[0] = byte(int8(v))
	default:
		b[0] = byte(int8(0))
	}
	return b
}

// Bytes32 bytes32
func Bytes32(input interface{}) []byte {
	switch v := input.(type) {
	case [32]byte:
		return RightPadBytes(v[:], 32)
	case []byte:
		return RightPadBytes(v, 32)
	case string:
		str := fmt.Sprintf("%x", v)
		hexb, _ := hex.DecodeString(str)
		return RightPadBytes(hexb, 32)
	default:
		return RightPadBytes([]byte(""), 32)
	}
}

// String string
func String(input interface{}) []byte {
	switch v := input.(type) {
	case []byte:
		return v
	case string:
		return []byte(v)
	default:
		return []byte("")
	}
}

// Bool bool
func Bool(input interface{}) []byte {
	switch v := input.(type) {
	case bool:
		if v {
			return []byte{0x1}
		}
		return []byte{0x0}
	default:
		return []byte{0x0}
	}
}

// ConcatByteSlices concat byte slices
func ConcatByteSlices(arrays ...[]byte) []byte {
	var result []byte

	for _, b := range arrays {
		result = append(result, b...)
	}

	return result
}
func U256(n *big.Int) []byte {
	return PaddedBigBytes(mathU256(n), 32)
}

func SoliditySHA3(data ...[]byte) []byte {
	var result []byte

	hash := sha3.NewKeccak256()
	bs := ConcatByteSlices(data...)

	hash.Write(bs)
	result = hash.Sum(result)

	return result
}

// Uint256Array uint256 array
func Uint256Array(input interface{}) []byte {
	var values []byte
	s := reflect.ValueOf(input)
	for i := 0; i < s.Len(); i++ {
		val := s.Index(i).Interface()
		result := LeftPadBytes(Uint256(val), 32)
		values = append(values, result...)
	}
	return values
}

// LeftPadBytes zero-pads slice to the left up to length l.
func LeftPadBytes(slice []byte, l int) []byte {
	if l <= len(slice) {
		return slice
	}

	padded := make([]byte, l)
	copy(padded[l-len(slice):], slice)

	return padded
}

// HexToAddress returns Address with byte values of s.
// If s is larger than len(h), s will be cropped from the left.
func HexToAddress(s string) Address { return BytesToAddress(FromHex(s)) }

// has0xPrefix validates str begins with '0x' or '0X'.
func has0xPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

// FromHex returns the bytes represented by the hexadecimal string s.
// s may be prefixed with "0x".
func FromHex(s string) []byte {
	if has0xPrefix(s) {
		s = s[2:]
	}
	if len(s)%2 == 1 {
		s = "0" + s
	}
	return Hex2Bytes(s)
}

// Hex2Bytes returns the bytes represented by the hexadecimal string str.
func Hex2Bytes(str string) []byte {
	h, _ := hex.DecodeString(str)
	return h
}

// RightPadBytes zero-pads slice to the right up to length l.
func RightPadBytes(slice []byte, l int) []byte {
	if l <= len(slice) {
		return slice
	}

	padded := make([]byte, l)
	copy(padded, slice)

	return padded
}

func PaddedBigBytes(bigint *big.Int, n int) []byte {
	if bigint.BitLen()/8 >= n {
		return bigint.Bytes()
	}
	ret := make([]byte, n)
	ReadBits(bigint, ret)
	return ret
}

const (
	// number of bits in a big.Word
	wordBits = 32 << (uint64(^big.Word(0)) >> 63)
	// number of bytes in a big.Word
	wordBytes = wordBits / 8
)

var (
	tt255   = BigPow(2, 255)
	tt256   = BigPow(2, 256)
	tt256m1 = new(big.Int).Sub(tt256, big.NewInt(1))
)

// BigPow returns a ** b as a big integer.
func BigPow(a, b int64) *big.Int {
	r := big.NewInt(a)
	return r.Exp(r, big.NewInt(b), nil)
}

// ReadBits encodes the absolute value of bigint as big-endian bytes. Callers must ensure
// that buf has enough space. If buf is too short the result will be incomplete.
func ReadBits(bigint *big.Int, buf []byte) {
	i := len(buf)
	for _, d := range bigint.Bits() {
		for j := 0; j < wordBytes && i > 0; j++ {
			i--
			buf[i] = byte(d)
			d >>= 8
		}
	}
}

// U256 encodes as a 256 bit two's complement number. This operation is destructive.
func mathU256(x *big.Int) *big.Int {
	return x.And(x, tt256m1)
}
