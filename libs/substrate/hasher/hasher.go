package hasher

import (
	"github.com/itering/subscan/util/twox"
	"golang.org/x/crypto/blake2b"
)

// HashByCryptoName
func HashByCryptoName(p []byte, hasher string) []byte {
	h := p
	switch hasher {
	case "Blake2_128":
		checksum, _ := blake2b.New(16, []byte{})
		checksum.Write(p)
		h = checksum.Sum(nil)
	case "Blake2_256":
		checksum, _ := blake2b.New256([]byte{})
		checksum.Write(p)
		h = checksum.Sum(nil)
	case "Twox128":
		p := twox.NewXXHash128(p)
		h = p[:]
	case "Twox256":
		p := twox.NewXXHash128(p)
		h = p[:]
	case "Twox64Concat":
		p := twox.TwoX64Concat(p)
		h = p[:]
	case "Identity":
		h = p[:]
	case "Blake2_128Concat":
		checksum, _ := blake2b.New(16, []byte{})
		checksum.Write(p)
		h = checksum.Sum(nil)
		h = append(h, p...)
	default:
		p := twox.TwoX64Concat(p)
		h = p[:]
	}
	return h
}
