package service

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/util"
	"github.com/itering/subscan/util/address"
	"github.com/itering/substrate-api-rpc/metadata"
)

// Ping ping the resource.
func (s *Service) Ping(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, s.dao.Ping(ctx)
}

func (s *Service) Metadata(ctx context.Context) (map[string]string, error) {
	m, err := s.dao.GetMetadata(ctx)
	if err != nil {
		return nil, err
	}
	m["networkNode"] = util.NetworkNode
	m["balanceAccuracy"] = util.BalanceAccuracy
	m["addressType"] = util.AddressType
	return m, err
}

func (s *Service) GetBlocksSampleByNums(ctx context.Context, page, row int) []model.SampleBlockJson {
	var blockJson []model.SampleBlockJson
	blocks := s.dao.GetBlockList(ctx, page, row)
	for _, block := range blocks {
		bj := s.BlockAsSampleJson(&block)
		blockJson = append(blockJson, *bj)
	}
	return blockJson
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

func (s *Service) GetBestBlockNum(c context.Context) (uint64, error) {
	return s.dao.GetBestBlockNum(c)
}

func (s *Service) GetExtrinsicList(ctx context.Context, page, row int, query ...string) ([]*model.ChainExtrinsicJson, int) {
	list, count := s.dao.GetExtrinsicList(ctx, page, row, "desc", query...)
	var ejs []*model.ChainExtrinsicJson
	for _, extrinsic := range list {
		ejs = append(ejs, s.dao.ExtrinsicsAsJson(&extrinsic))
	}
	return ejs, count
}

func (s *Service) GetExtrinsicByIndex(ctx context.Context, index string) *model.ExtrinsicDetail {
	return s.dao.GetExtrinsicsDetailByIndex(ctx, index)
}

func (s *Service) GetExtrinsicDetailByHash(ctx context.Context, hash string) *model.ExtrinsicDetail {
	return s.dao.GetExtrinsicsDetailByHash(ctx, hash)
}

func (s *Service) EventsList(ctx context.Context, page, row int, where ...string) ([]model.ChainEventJson, int) {
	var (
		result    []model.ChainEventJson
		blockNums []uint
	)

	list, count := s.dao.GetEventList(ctx, page, row, "desc", where...)
	for _, event := range list {
		blockNums = append(blockNums, event.BlockNum)
	}
	blockMap := s.dao.BlocksReverseByNum(blockNums)

	for _, event := range list {
		ej := model.ChainEventJson{
			ExtrinsicIdx: event.ExtrinsicIdx,
			EventIndex:   event.EventIndex,
			BlockNum:     event.BlockNum,
			ModuleId:     event.ModuleId,
			EventId:      event.EventId,
			Params:       util.ToString(event.Params),
			EventIdx:     event.EventIdx,
		}
		if block, ok := blockMap[event.BlockNum]; ok {
			ej.BlockTimestamp = block.BlockTimestamp
		}
		result = append(result, ej)
	}
	return result, count
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
