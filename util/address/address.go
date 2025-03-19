package address

import (
	"github.com/itering/subscan/util"
	"github.com/itering/subscan/util/ss58"
	"regexp"
	"strings"
)

var (
	ethAddressRegex       = regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`)
	substrateAddressRegex = regexp.MustCompile(`^[0-9a-fA-F]{64}$`)
)

func VerifyEthereumAddress(accountId string) bool {
	return ethAddressRegex.MatchString(util.AddHex(accountId))
}

func SS58Address(address string) string {
	return ss58.Encode(address, util.StringToInt(util.AddressType))
}

func Encode(address string) string {
	if VerifyEthereumAddress(address) {
		return util.AddHex(address)
	}
	return SS58Address(util.TrimHex(address))
}

func VerifySubstrateAddress(accountId string) bool {
	return substrateAddressRegex.MatchString(util.TrimHex(accountId))
}

func Decode(address string) string {
	if VerifyEthereumAddress(address) {
		return util.AddHex(strings.ToLower(address))
	}
	return ss58.Decode(address)
}

func Format(accountID string) string {
	if VerifySubstrateAddress(accountID) {
		return util.TrimHex(accountID)
	}
	// Ethereum Address
	if VerifyEthereumAddress(accountID) {
		return util.AddHex(strings.ToLower(accountID))
	}
	return ""
}
