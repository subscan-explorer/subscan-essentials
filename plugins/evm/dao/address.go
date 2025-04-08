package dao

import (
	"github.com/itering/subscan/util"
)

func RemoveAddressPadded(address string) string {
	address = util.TrimHex(address)
	if len(address) < 40 {
		return NullAddress
	}
	return util.AddHex(address[len(address)-40:])
}
