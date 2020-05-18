package substrate

import (
	"golang.org/x/crypto/blake2b"
	"subscan-end/utiles/twox"
)

func hashBytesByHasher(p []byte, hasher string) []byte {
	h := p
	switch hasher {
	case "Blake2_128":
		checksum, _ := blake2b.New(128, []byte{})
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
	case "Twox128Concat":
	}
	return h
}
