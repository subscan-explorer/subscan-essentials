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

// SS58Address converts the address to SS58 format
func SS58Address(address string) string {
	return ss58.Encode(address, util.StringToInt(util.AddressType))
}

// VerifyEthereumAddress checks if the address is a valid Ethereum address
func VerifyEthereumAddress(accountId string) bool {
	return ethAddressRegex.MatchString(util.AddHex(accountId))
}

// Encode converts the address to Ethereum or SS58 format, depending on the address type Ethereum or Substrate
func Encode(address string) string {
	if VerifyEthereumAddress(address) {
		return util.AddHex(address)
	}
	return SS58Address(util.TrimHex(address))
}

// VerifySubstrateAddress checks if the address is a valid Substrate address
func VerifySubstrateAddress(accountId string) bool {
	return substrateAddressRegex.MatchString(util.TrimHex(accountId))
}

// Decode converts the address to Substrate public key or Ethereum format, depending on the address type Ethereum or Substrate
func Decode(address string) string {
	if VerifyEthereumAddress(address) {
		return util.AddHex(strings.ToLower(address))
	}
	return ss58.Decode(address)
}

// Format accountId to db format, ethereum address to lowercase and add 0x
// substrate address to lowercase and remove 0x
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
