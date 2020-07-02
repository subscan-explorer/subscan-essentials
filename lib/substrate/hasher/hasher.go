package hasher

import (
	"github.com/itering/subscan/util/twox"
	"golang.org/x/crypto/blake2b"
)

// HashByCryptoName
func HashByCryptoName(p []byte, hasher string) []byte {
	switch hasher {
	case "Blake2_128":
		checksum, _ := blake2b.New(16, []byte{})
		_, _ = checksum.Write(p)
		p = checksum.Sum(nil)
	case "Blake2_256":
		checksum, _ := blake2b.New256([]byte{})
		_, _ = checksum.Write(p)
		p = checksum.Sum(nil)
	case "Twox128":
		h := twox.NewXXHash128(p)
		p = h[:]
	case "Twox256":
		h := twox.NewXXHash128(p)
		p = h[:]
	case "Twox64Concat":
		h := twox.To64Concat(p)
		p = h[:]
	case "Identity":
		p = p[:]
	case "Blake2_128Concat":
		checksum, _ := blake2b.New(16, []byte{})
		_, _ = checksum.Write(p)
		h := checksum.Sum(nil)
		p = append(h, p...)
	default:
		h := twox.To64Concat(p)
		p = h[:]
	}
	return p
}
