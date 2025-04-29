package dao

import (
	"context"
	subscan_plugin "github.com/itering/subscan-plugin"
	"github.com/itering/subscan-plugin/storage"
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/share/token"
	"github.com/itering/subscan/util"
)

type Storage struct {
	Dao  storage.Dao
	Pool subscan_plugin.RedisPool
}

func EmitEvent(ctx context.Context, d *Storage, event *storage.Event, block *storage.Block) error {
	var paramEvent []storage.EventParam
	_ = util.UnmarshalAny(&paramEvent, event.Params)

	switch event.EventId {
	// [accountId, balance]
	case "Endowed", "Reserved", "Unreserved", "Deposit", "Minted", "Issued", "Locked", "Unlocked", "Withdraw":
		return RefreshAccount(ctx, d, model.CheckoutParamValueAddress(paramEvent[0].Value))
		// ["AccountId","AccountId","Balance"]
	case "Transfer":
		from := model.CheckoutParamValueAddress(paramEvent[0].Value)
		to := model.CheckoutParamValueAddress(paramEvent[1].Value)
		balance := util.DecimalFromInterface(paramEvent[2].Value)
		t := token.GetDefaultToken()
		return CreateTransfer(ctx, d, &Transfer{
			Id:             event.Id,
			Sender:         from,
			Receiver:       to,
			Amount:         balance,
			BlockTimestamp: int64(block.BlockTimestamp),
			Symbol:         t.Symbol,
			TokenId:        t.TokenId,
		})
	}
	return nil
}
