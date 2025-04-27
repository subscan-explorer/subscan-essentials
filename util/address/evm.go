package address

import (
	"github.com/itering/subscan/util"

	"golang.org/x/crypto/blake2b"
)

func SS58AddressToEvm(ss58Address string) string {
	ss58Address = util.TrimHex(ss58Address)
	if len(ss58Address) != 64 {
		return ""
	}
	return util.AddHex(ss58Address[0:40])
}

func EvmToSS58Address(evmAddress string) string {
	if !VerifyEthereumAddress(evmAddress) {
		return ""
	}
	checksum, _ := blake2b.New(32, []byte{})
	_, _ = checksum.Write(append([]byte("evm:"), util.HexToBytes(evmAddress)...))
	return util.BytesToHex(checksum.Sum(nil))
}
