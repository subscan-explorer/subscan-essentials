package address

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SS58AddressToEvm(t *testing.T) {
	assert.Equal(t, "0xa7839fbfca6da129ff9e2ff521115b7eb4213b21", SS58AddressToEvm("a7839fbfca6da129ff9e2ff521115b7eb4213b215086fc3416c7a340e944cc49"))
}

func Test_EvmToSS58Address(t *testing.T) {
	assert.Equal(t, "a7839fbfca6da129ff9e2ff521115b7eb4213b215086fc3416c7a340e944cc49", EvmToSS58Address("0xd43593c715fdd31c61141abd04a99fd6822c8558"))
}
