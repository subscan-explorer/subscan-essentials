package address

import (
	"strings"

	"github.com/itering/subscan/util"
	"github.com/itering/subscan/util/ss58"
)

func SS58AddressFromHex(address string) SS58Address {
	address = strings.TrimPrefix(address, "0x")
	addressType := util.StringToInt(util.AddressType)

	return SS58Address(ss58.Encode(address, addressType))
}

type SS58Address string

func (a SS58Address) String() string {
	return string(a)
}
