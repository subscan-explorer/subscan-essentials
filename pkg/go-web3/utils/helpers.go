package utils

import (
	"fmt"
	"golang.org/x/crypto/sha3"
	"hash"
	"math/big"
)

func IntToHex(n *big.Int) string {
	return fmt.Sprintf("0x%x", n)
}

type KeccakState interface {
	hash.Hash
	Read([]byte) (int, error)
}

func Keccak256(data ...[]byte) []byte {
	b := make([]byte, 32)
	d := sha3.NewLegacyKeccak256().(KeccakState)
	for _, b := range data {
		d.Write(b)
	}
	d.Read(b)
	return b
}
