package web3

import (
	"github.com/itering/subscan/pkg/go-web3"
	"github.com/itering/subscan/pkg/go-web3/providers"
	"github.com/itering/subscan/util"
	"os"
)

var (
	RPC *web3.Web3
)

// assethub westend https://westend-asset-hub-eth-rpc.polkadot.io
func init() {
	if EthRpc := os.Getenv("ETH_RPC"); EthRpc != "" {
		RPC = web3.NewWeb3(providers.NewHTTPProvider(EthRpc, 60, false))
		return
	}
	RPC = web3.NewWeb3(providers.NewWebSocketProvider(util.WSEndPoint))
}
