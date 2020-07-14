package address

import (
	"github.com/itering/subscan/util"
	"github.com/itering/subscan/util/ss58"
)

func SS58Address(address string) string {
	return ss58.Encode(address, util.StringToInt(util.AddressType))
}
