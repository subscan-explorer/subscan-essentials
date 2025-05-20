package abi

import (
	"encoding/hex"
	"fmt"
	"strings"

	solsha3 "github.com/itering/subscan/pkg/go-solidity-sha3"
)

func EncodingMethod(content string) string {
	hash := solsha3.SoliditySHA3(
		solsha3.String(content),
	)
	return hex.EncodeToString(hash)
}

func DecodeStaticType(str string) string {
	if len(str) == 0 {
		return str
	}
	return strings.TrimLeft(strings.TrimPrefix(str, "0x"), "0")
}

func DecodeAddress(str string) string {
	address := DecodeStaticType(str)
	if address == "" {
		return address
	}
	return fmt.Sprintf("%040s", address)
}
