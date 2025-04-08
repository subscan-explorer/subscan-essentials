package delegateProxy

import (
	"context"
	"github.com/itering/subscan/pkg/go-web3"
	"github.com/itering/subscan/pkg/go-web3/providers"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDelegateProxy(t *testing.T) {
	ctx := context.Background()
	proxy897 := Init897(web3.NewWeb3(providers.NewHTTPProvider("https://eth-mainnet.g.alchemy.com/v2/TfPV2XK_Xw02-1jNjb-jsMKju1R14lm9", 60, false)), "0x14a4123da9ad21b2215dc0ab6984ec1e89842c6d")
	Implementation, err := proxy897.Implementation(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "0x85520f613021e5db2afb40cb91ef066ccf212111", Implementation)
}

func TestEIP1967_Implementation(t *testing.T) {
	ctx := context.Background()
	proxy1967 := Init1967(
		web3.NewWeb3(providers.NewHTTPProvider("https://moonbeam-rpc.n.dwellir.com", 60, false)),
		"0x25442adf37379be90ed1f7fccd9c9417b10aa4dc")
	Implementation, err := proxy1967.Implementation(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "0x40f1eca9c82200428704aa555da1009ad4beb2e2", Implementation)

	// no implementation address is 0x0000000000000000000000000000000000000000
	proxy1967.Contract.TransParam.To = "0x22b1a40e3178fe7c7109efcc247c5bb2b34abe32"
	Implementation, err = proxy1967.Implementation(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "0x0000000000000000000000000000000000000000", Implementation)

}
