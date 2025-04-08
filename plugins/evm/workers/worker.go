package workers

import (
	"context"
	"github.com/itering/subscan/plugins/evm/dao"
	"github.com/itering/subscan/plugins/evm/feature/erc1155"
	"github.com/itering/subscan/util"
)

func Emit(ctx context.Context, queue, class string, raw interface{}) error {
	switch queue {
	case dao.Eip20Token, dao.Eip721Token:
		switch class {
		case "transfer":
			var receipt dao.TransactionReceipt
			util.Logger().Error(util.UnmarshalAny(&receipt, raw))
			return receipt.ProcessTokenTransfer(ctx, queue)

		case "balance":
			// [contract, address]
			var contractAddress []string
			util.Logger().Error(util.UnmarshalAny(&contractAddress, raw))
			return dao.RefreshHolder(ctx, contractAddress[0], contractAddress[1], queue)

		case "holder":
			// [contract, tokenId]
			var args []string
			util.Logger().Error(util.UnmarshalAny(&args, raw))
			if token := dao.GetTokenByContract(ctx, args[0]); token != nil {
				return token.RefreshErc721Holders(ctx, args[1])
			}
		}
	case dao.Eip1155Token:
		switch class {
		case "balance":
			// [contract, address, tokenIds []string ]
			var contractAddress []string
			util.Logger().Error(util.UnmarshalAny(&contractAddress, raw))
			if contractAddress[1] == dao.NullAddress {
				return nil
			}
			return dao.RefreshErc1155Holder(ctx, contractAddress[0], contractAddress[1], contractAddress[2])
		default:
			var receipt dao.TransactionReceipt
			util.Logger().Error(util.UnmarshalAny(&receipt, raw))
			switch util.TrimHex(receipt.MethodHash) {
			case erc1155.EventTransferBatch, erc1155.EventTransferSingle, erc1155.EventURI:
				return receipt.ProcessErc1155(ctx)
			}
		}
	}
	return nil
}
