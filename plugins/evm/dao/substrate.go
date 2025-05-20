package dao

import (
	"context"
	"fmt"
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/pkg/go-web3/dto"
	"github.com/itering/subscan/plugins/evm/feature"
	"github.com/itering/subscan/util"
	"github.com/itering/substrate-api-rpc/metadata"
	"strings"
)

func findOutSubstrateExecutedEvent(ctx context.Context, blockNum uint, _ *dto.Block) (hash2ExtrinsicIndex map[string]string, err error) {
	supportModule := metadata.SupportModule()

	hash2ExtrinsicIndex = make(map[string]string)

	switch {
	case util.StringInSlice("Ethereum", supportModule):
		extrinsicIndex2Events := GetEventsByBlockNum(ctx, blockNum, model.Where("module_id = ? AND event_id = ?", "ethereum", "Executed"))
		for _, event := range extrinsicIndex2Events {
			if strings.EqualFold(fmt.Sprintf("%s.%s", event.ModuleId, event.EventId), "ethereum.Executed") {
				// [from, transaction_hash, ExitReason]
				if len(event.Params) == 3 {
					hash2ExtrinsicIndex[util.ToString(event.Params[1].Value)] = event.ExtrinsicIndex
					// [from, to/contract_address, transaction_hash, ExitReason]
				} else if len(event.Params) >= 4 {
					hash2ExtrinsicIndex[util.ToString(event.Params[2].Value)] = event.ExtrinsicIndex
				}
			}
		}
	case util.StringInSlice("Revive", supportModule):
		ethTransacts := GetExtrinsicsByBlockNum(ctx, blockNum, model.Where("call_module = ? AND call_module_function = ?", "revive", "eth_transact"))
		for _, ethTransact := range ethTransacts {
			if len(ethTransact.Params) == 0 {
				continue
			}
			txHash := feature.CalHashByTxRaw(ctx, ethTransact.Params[0].Value.(string))
			if txHash == "" {
				continue
			}
			hash2ExtrinsicIndex[txHash] = ethTransact.ExtrinsicIndex
		}
	}
	return hash2ExtrinsicIndex, nil
}

func GetEventsByBlockNum(c context.Context, blockNum uint, options ...model.Option) []model.ChainEvent {
	var Event []model.ChainEvent
	query := sg.db.WithContext(c).Scopes(model.TableNameFunc(&model.ChainEvent{BlockNum: blockNum})).
		Where("block_num = ?", blockNum)
	query = query.Scopes(options...)
	query = query.Find(&Event)
	if query == nil || query.Error != nil {
		return nil
	}
	return Event
}

func GetExtrinsicsByBlockNum(c context.Context, blockNum uint, opts ...model.Option) []model.ChainExtrinsic {
	var extrinsics []model.ChainExtrinsic
	q := sg.db.WithContext(c).Scopes(model.TableNameFunc(&model.ChainExtrinsic{BlockNum: blockNum})).Scopes(opts...).
		Where("block_num = ?", blockNum).Order("id asc").Find(&extrinsics)
	if q == nil || q.Error != nil {
		return nil
	}
	return extrinsics
}
