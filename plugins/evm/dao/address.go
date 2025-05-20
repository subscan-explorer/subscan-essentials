package dao

import (
	"context"
	"github.com/huandu/xstrings"
	"github.com/itering/subscan/util"
	"github.com/itering/subscan/util/address"
	"github.com/itering/substrate-api-rpc/metadata"
)

func RemoveAddressPadded(address string) string {
	address = util.TrimHex(address)
	if len(address) < 40 {
		return NullAddress
	}
	return util.AddHex(address[len(address)-40:])
}

func h160ToAccountIdByNetwork(_ context.Context, h160 string, _ string) string {
	if address.VerifyEthereumAddress(h160) {
		h160 = address.Format(h160)
		switch {
		case util.IsEvmChain:
			return h160
		case metadata.Latest(nil) != nil && util.StringInSlice("Revive", metadata.SupportModule()):
			return reviveAccount(h160)
		default:
			return address.EvmToSS58Address(h160)
		}
	}
	return ""
}

func reviveAccount(h160 string) string {
	return xstrings.LeftJustify(util.TrimHex(h160), 64, "e")
}
