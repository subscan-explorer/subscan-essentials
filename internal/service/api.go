package service

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/itering/subscan/configs"
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/plugins"
	"github.com/itering/subscan/util"
	"github.com/itering/subscan/util/address"
	"github.com/itering/substrate-api-rpc/metadata"
)

// Ping ping the resource.
func (s *Service) Ping(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, s.dao.Ping(ctx)
}

func (s *Service) Metadata(ctx context.Context) (map[string]interface{}, error) {
	meta, err := s.dao.GetMetadata(ctx)
	if err != nil {
		return nil, err
	}
	m := make(map[string]interface{})
	m["networkNode"] = util.NetworkNode
	m["balanceAccuracy"] = util.BalanceAccuracy
	m["addressType"] = util.AddressType

	evm, ok := plugins.RegisteredPlugins["evm"]
	if ok && evm.Enable() {
		m["enable_evm"] = configs.Boot.UI.EnableEvm
	}
	m["enable_substrate"] = configs.Boot.UI.EnableSubstrate
	for k, v := range meta {
		m[k] = v
	}
	return m, err
}

func (s *Service) GetBlocksSampleCursor(ctx context.Context, limit int, before, after uint) ([]model.SampleBlockJson, CursorPage) {
	var blockJson []model.SampleBlockJson
	blocks, hasPrev, hasNext := s.dao.GetBlockListCursor(ctx, limit, before, after)
	for _, block := range blocks {
		bj := s.BlockAsSampleJson(&block)
		blockJson = append(blockJson, *bj)
	}
	var start, end *uint
	if len(blocks) > 0 {
		startBlock := blocks[0].BlockNum
		endBlock := blocks[len(blocks)-1].BlockNum
		start = &startBlock
		end = &endBlock
	}
	return blockJson, CursorPage{StartCursor: start, EndCursor: end, HasNextPage: hasNext, HasPreviousPage: hasPrev}
}

func (s *Service) BlockAsSampleJson(block *model.ChainBlock) *model.SampleBlockJson {
	b := model.SampleBlockJson{
		BlockNum:        block.BlockNum,
		BlockTimestamp:  block.BlockTimestamp,
		Hash:            block.Hash,
		EventCount:      block.EventCount,
		ExtrinsicsCount: block.ExtrinsicsCount,
		Validator:       address.Encode(block.Validator),
		Finalized:       block.Finalized,
	}
	return &b
}

func (s *Service) GetBlockByNum(ctx context.Context, num uint) *model.ChainBlockJson {
	block := s.dao.GetBlockByNum(ctx, num)
	if block == nil {
		return nil
	}
	return s.dao.BlockAsJson(ctx, block)
}

func (s *Service) GetBlockByHash(ctx context.Context, hash string) *model.ChainBlock {
	block := s.dao.GetBlockByHash(ctx, hash)
	if block == nil {
		return nil
	}
	return block
}

func (s *Service) GetFinalizedBlock(c context.Context) (uint64, error) {
	return s.dao.GetFinalizedBlockNum(c)
}

type CursorPage struct {
	StartCursor     *uint `json:"start_cursor,omitempty"`
	EndCursor       *uint `json:"end_cursor,omitempty"`
	HasNextPage     bool  `json:"has_next_page"`
	HasPreviousPage bool  `json:"has_previous_page"`
}

func (s *Service) GetExtrinsicList(ctx context.Context, limit int, fixedTableIndex int, beforeId, afterId uint, accountId string, query ...model.Option) ([]*model.ChainExtrinsicJson, CursorPage) {
	list, hasPrev, hasNext := s.dao.GetExtrinsicListCursor(ctx, limit, fixedTableIndex, beforeId, afterId, accountId, query...)
	var ejs []*model.ChainExtrinsicJson
	for _, extrinsic := range list {
		ejs = append(ejs, s.dao.ExtrinsicsAsJson(&extrinsic))
	}
	var start, end *uint
	if len(list) > 0 {
		start = &list[0].ID
		end = &list[len(list)-1].ID
	}
	return ejs, CursorPage{StartCursor: start, EndCursor: end, HasNextPage: hasNext, HasPreviousPage: hasPrev}
}

func (s *Service) GetExtrinsicByIndex(ctx context.Context, index string) *model.ExtrinsicDetail {
	return s.dao.GetExtrinsicsDetailByIndex(ctx, index)
}

func (s *Service) GetExtrinsicDetailByHash(ctx context.Context, hash string) *model.ExtrinsicDetail {
	return s.dao.GetExtrinsicsDetailByHash(ctx, hash)
}

func (s *Service) EventsList(ctx context.Context, limit int, fixedTableIndex int, beforeId uint, afterId uint, where ...model.Option) ([]model.ChainEventJson, CursorPage) {
	var (
		result    []model.ChainEventJson
		blockNums []uint
	)

	list, hasPrev, hasNext := s.dao.GetEventListCursor(ctx, limit, "desc", fixedTableIndex, beforeId, afterId, where...)
	for _, event := range list {
		blockNums = append(blockNums, event.BlockNum)
	}
	blockMap := s.dao.BlocksReverseByNum(blockNums)

	for _, event := range list {
		ej := model.ChainEventJson{
			Id:             event.ID,
			ExtrinsicIndex: event.ExtrinsicIndex,
			BlockNum:       event.BlockNum,
			ModuleId:       event.ModuleId,
			EventId:        event.EventId,
			Params:         event.Params,
			EventIdx:       event.EventIdx,
			EventIndex:     fmt.Sprintf("%d-%d", event.BlockNum, event.EventIdx),
			Phase:          event.Phase,
		}
		if block, ok := blockMap[event.BlockNum]; ok {
			ej.BlockTimestamp = block.BlockTimestamp
		}
		result = append(result, ej)
	}
	var start, end *uint
	if len(list) > 0 {
		start = &list[0].ID
		end = &list[len(list)-1].ID
	}
	return result, CursorPage{StartCursor: start, EndCursor: end, HasNextPage: hasNext, HasPreviousPage: hasPrev}
}

func (s *Service) EventById(ctx context.Context, eventIndex string) *model.ChainEventJson {
	event := s.dao.GetEventByIdx(ctx, eventIndex)
	if event == nil {
		return nil
	}
	ej := model.ChainEventJson{
		Id:             event.ID,
		ExtrinsicIndex: event.ExtrinsicIndex,
		BlockNum:       event.BlockNum,
		ModuleId:       event.ModuleId,
		EventId:        event.EventId,
		Params:         event.Params,
		EventIdx:       event.EventIdx,
		EventIndex:     fmt.Sprintf("%d-%d", event.BlockNum, event.EventIdx),
		Phase:          event.Phase,
	}
	block := s.dao.GetBlockByNum(ctx, event.BlockNum)
	if block != nil {
		ej.BlockTimestamp = block.BlockTimestamp
	}
	return &ej
}

func (s *Service) GetExtrinsicByHash(ctx context.Context, hash string) *model.ChainExtrinsic {
	return s.dao.GetExtrinsicsByHash(ctx, hash)
}

func (s *Service) SubstrateRuntimeList() []model.RuntimeVersion {
	return s.dao.RuntimeVersionList()
}

func (s *Service) SubstrateRuntimeInfo(spec int) *metadata.Instant {
	if metadataInstant, ok := metadata.RuntimeMetadata[spec]; ok {
		return metadataInstant
	}
	runtime := metadata.Process(s.dao.RuntimeVersionRaw(spec))
	if runtime == nil {
		return metadata.Latest(nil)
	}
	return runtime
}

func (s *Service) LogsList(ctx context.Context, blockNum uint) []model.ChainLogJson {
	return s.dao.GetLogByBlockNum(ctx, blockNum)
}
