package dao

import (
	"context"
	"github.com/itering/subscan/plugins/evm/abi"
	"github.com/itering/subscan/plugins/evm/feature/delegateProxy"
	"github.com/itering/subscan/plugins/evm/feature/erc20"
	"github.com/itering/subscan/plugins/evm/feature/erc721"
	"github.com/itering/subscan/share/web3"
	"github.com/itering/subscan/util"
	"strings"
)

// EventProcess Process TransactionReceipt event log
func (t TransactionReceipt) EventProcess(ctx context.Context) error {
	methodSha3 := util.TrimHex(t.MethodHash)

	switch methodSha3 {

	// check erc20 or erc721
	case erc20.EventTransfer, erc721.EventTransfer:
		// 符合 transfer event 默认归类到ERC 20, 需要进一步确认是否是ERC 721
		topics := strings.Split(t.Topics, ",")
		tokenId := t.Data
		if len(topics) == 4 {
			tokenId = topics[3]
		}
		if token := GetTokenByContract(ctx, t.Address); token == nil || token.Category == "" {
			if erc721.BaseVerify(ctx, web3.RPC, t.Address, BillionAddress(ctx), util.TrimHex(tokenId)) {
				_ = Publish(Eip721Token, "transfer", t)
				return nil
			}
			_ = Publish(Eip20Token, "transfer", t)
		} else {
			_ = Publish(token.Category, "transfer", t)
		}
	// deposit or withdraw
	case erc20.EventDeposit, erc20.EventWithdraw:
		log2 := strings.Split(t.Topics, ",")
		if len(log2) != 2 {
			return nil
		}
		if token := GetTokenByContract(ctx, t.Address); token == nil {
			return nil
		}
		address := util.AddHex(abi.DecodeAddress(log2[1]))
		_ = Publish("erc20", "balance", []string{t.Address, address})

	// proxy
	case delegateProxy.EventUpgraded:
		if topics := strings.Split(t.Topics, ","); len(topics) > 1 {
			setContractProxyImplementation(ctx, t.Address, util.AddHex(abi.DecodeAddress(topics[1])))
		}

		// case erc1155.EventTransferBatch, erc1155.EventTransferSingle, erc1155.EventURI:
		// 	if token := GetTokenByContract(ctx, t.Address); token == nil {
		// 		erc1155token := erc1155.Init(web3.RPC, t.Address)
		// 		if result, _ := erc1155token.SupportsInterface(ctx); result {
		// 			return Publish(Eip1155Token, Eip1155Token, t)
		// 		}
		// 	} else if token.Category == Eip1155Token {
		// 		_ = Publish(token.Category, Eip1155Token, t)
		// 	}
	}
	return nil
}
