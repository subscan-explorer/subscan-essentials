package ss58

import (
	"bytes"
	"encoding/binary"
	"github.com/freehere107/go-scale-codec/types"
	"github.com/itering/subscan/internal/util"
	"github.com/itering/subscan/internal/util/base58"
	"golang.org/x/crypto/blake2b"
)

func Decode(address string, addressType int) string {
	checksumPrefix := []byte("SS58PRE")
	ss58Format := base58.Decode(address)
	if len(ss58Format) == 0 || ss58Format[0] != byte(addressType) {
		return ""
	}
	var checksumLength int
	if util.IntInSlice(len(ss58Format), []int{3, 4, 6, 10}) {
		checksumLength = 1
	} else if util.IntInSlice(len(ss58Format), []int{5, 7, 11, 35}) {
		checksumLength = 2
	} else if util.IntInSlice(len(ss58Format), []int{8, 12}) {
		checksumLength = 3
	} else if util.IntInSlice(len(ss58Format), []int{9, 13}) {
		checksumLength = 4
	} else if util.IntInSlice(len(ss58Format), []int{14}) {
		checksumLength = 5
	} else if util.IntInSlice(len(ss58Format), []int{15}) {
		checksumLength = 6
	} else if util.IntInSlice(len(ss58Format), []int{16}) {
		checksumLength = 7
	} else if util.IntInSlice(len(ss58Format), []int{17}) {
		checksumLength = 8
	} else {
		return ""
	}
	bss := ss58Format[0 : len(ss58Format)-checksumLength]
	checksum, _ := blake2b.New(64, []byte{})
	w := append(checksumPrefix[:], bss[:]...)
	_, err := checksum.Write(w)
	if err != nil {
		return ""
	}

	h := checksum.Sum(nil)
	if util.BytesToHex(h[0:checksumLength]) != util.BytesToHex(ss58Format[len(ss58Format)-checksumLength:]) {
		return ""
	}
	return util.BytesToHex(ss58Format[1 : len(ss58Format)-checksumLength])
}

func Encode(address string, addressType int) string {
	checksumPrefix := []byte("SS58PRE")
	addressBytes := util.HexToBytes(address)
	var checksumLength int
	if len(addressBytes) == 32 {
		checksumLength = 2
	} else if util.IntInSlice(len(addressBytes), []int{1, 2, 4, 8}) {
		checksumLength = 1
	} else {
		return ""
	}
	addressFormat := append([]byte{byte(addressType)}[:], addressBytes[:]...)
	checksum, _ := blake2b.New(64, []byte{})
	w := append(checksumPrefix[:], addressFormat[:]...)
	_, err := checksum.Write(w)
	if err != nil {
		return ""
	}

	h := checksum.Sum(nil)
	b := append(addressFormat[:], h[:checksumLength][:]...)
	return base58.Encode(b)
}

func EncodeAccountIndex(accountIndex int64, addressType int) string {
	if accountIndex == -1 {
		return ""
	}
	buf := new(bytes.Buffer)
	var err error
	if accountIndex >= 0 && accountIndex <= 255 {
		err = binary.Write(buf, binary.LittleEndian, byte(accountIndex))
	} else if accountIndex >= 256 && accountIndex <= 65535 {
		err = binary.Write(buf, binary.LittleEndian, uint16(accountIndex))
	} else if accountIndex >= 65536 && accountIndex <= 4294967295 {
		err = binary.Write(buf, binary.LittleEndian, uint32(accountIndex))
	} else {
		err = binary.Write(buf, binary.LittleEndian, uint64(accountIndex))
	}
	if err != nil {
		return ""
	}
	return Encode(util.BytesToHex(buf.Bytes()), addressType)
}

func DecodeAccountIndex(accountIndex string, addressType int) int64 {
	if accountIndex == "" {
		return -1
	}
	index := Decode(accountIndex, addressType)
	data := types.ScaleBytes{Data: util.HexToBytes(index)}
	switch len(index) {
	case 2:
		byteInstant := types.U8{}
		byteInstant.Init(data, nil)
		byteInstant.Process()
		return int64(byteInstant.Value.(int))
	case 4:
		byteInstant := types.U16{}
		byteInstant.Init(data, nil)
		byteInstant.Process()
		return int64(byteInstant.Value.(uint16))
	case 8:
		byteInstant := types.U32{}
		byteInstant.Init(data, nil)
		byteInstant.Process()
		return int64(byteInstant.Value.(uint32))
	case 16:
		byteInstant := types.U64{}
		byteInstant.Init(data, nil)
		byteInstant.Process()
		return int64(byteInstant.Value.(uint64))
	}
	return -1
}
