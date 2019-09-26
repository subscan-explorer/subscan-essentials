package twox

import (
	"github.com/pierrec/xxHash/xxHash64"
	"math"
)

func NewXXHash(data []byte, bitLength int64) []byte {
	return newXXHash(data, uint(bitLength))
}

func NewXXHash64(data []byte) [8]byte {
	var hash [8]byte
	copy(hash[:], newXXHash(data, 64))
	return hash
}

func NewXXHash128(data []byte) [16]byte {
	var hash [16]byte
	copy(hash[:], newXXHash(data, 128))
	return hash
}

func NewXXHash256(data []byte) [32]byte {
	var hash [32]byte
	copy(hash[:], newXXHash(data, 256))
	return hash
}

func newXXHash(data []byte, bitLength uint) []byte {
	byteLength := int64(math.Ceil(float64(bitLength) / float64(8)))
	iterations := int64(math.Ceil(float64(bitLength) / float64(64)))
	var hash = make([]byte, byteLength)

	for seed := int64(0); seed < iterations; seed++ {
		digest := xxHash64.New(uint64(seed))
		_, _ = digest.Write(data)
		copy(hash[seed*8:], digest.Sum(nil))
	}

	return hash
}
