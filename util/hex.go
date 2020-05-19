package util

import (
	"encoding/hex"
	"fmt"
	"strings"
)

// Hex
func AddHex(s string) string {
	if s == "" {
		return ""
	}
	if strings.HasPrefix(s, "0x") {
		return s
	}
	return strings.ToLower("0x" + s)
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

func IntToHex(i interface{}) string {
	return fmt.Sprintf("%x", i)
}

func TrimHex(s string) string {
	return strings.TrimPrefix(s, "0x")
}

// Bool
func BoolFromInterface(i interface{}) bool {
	switch i := i.(type) {
	case string:
		return strings.ToLower(i) == "true"
	case bool:
		return i
	}
	return false
}
