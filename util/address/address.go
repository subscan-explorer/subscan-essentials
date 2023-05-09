package address

import (
	"strings"

	"github.com/itering/subscan/util/ss58"
)

func SS58Address(address string) string {
	address = strings.TrimPrefix(address, "0x")
	return ss58.Encode(address, 42)
}
