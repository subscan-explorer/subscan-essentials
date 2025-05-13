package web3

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/itering/subscan/pkg/go-web3"
	"github.com/itering/subscan/pkg/go-web3/providers"
	"github.com/itering/subscan/util"
	"os"
)

var (
	RPC      *web3.Web3
	CHAIN_ID int64
)

// assethub westend https://westend-asset-hub-eth-rpc.polkadot.io
func init() {
	if EthRpc := os.Getenv("ETH_RPC"); EthRpc != "" {
		RPC = web3.NewWeb3(providers.NewHTTPProvider(EthRpc, 60, false))
	} else {
		RPC = web3.NewWeb3(providers.NewWebSocketProvider(util.WSEndPoint))
	}
	chainId, err := RPC.Eth.GetChainId(context.TODO())
	if err != nil {
		log.Debugf("get chain id error: %v, maybe not a evm chain", err)
	} else {
		CHAIN_ID = chainId.Int64()
	}
}
