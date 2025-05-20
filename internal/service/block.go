package service

import (
	"context"
	"fmt"
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/util"
	"github.com/itering/substrate-api-rpc"
	"github.com/itering/substrate-api-rpc/hasher"
	smodel "github.com/itering/substrate-api-rpc/model"
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/itering/substrate-api-rpc/storage"
)

func (s *Service) CreateChainBlock(ctx context.Context, hash string, block *smodel.Block, event string, spec int) (err error) {
	var (
		decodeExtrinsics []map[string]interface{}
		decodeEvent      interface{}
		logs             []storage.DecoderLog
		validator        string
		codecErr         error
	)

	blockNum := util.StringToUInt(util.HexToNumStr(block.Header.Number))
	metadataInstant := s.getMetadataInstant(spec, hash)

	// Extrinsic
	decodeExtrinsics, err = substrate.DecodeExtrinsic(block.Extrinsics, metadataInstant, spec)
	if err != nil {
		codecErr = err
		util.Logger().Error(err)
	}
	// event
	decodeEvent, err = substrate.DecodeEvent(event, metadataInstant, spec)
	if err != nil {
		codecErr = err
		util.Logger().Error(err)
	}

	// log
	logs, err = substrate.DecodeLogDigest(block.Header.Digest.Logs)
	if err != nil {
		codecErr = err
		util.Logger().Error(err)
	}

	txn := s.dao.DbBegin()
	defer s.dao.DbRollback(txn)

	var events []model.ChainEvent
	_ = util.UnmarshalAny(&events, decodeEvent)

	eventMap := s.checkoutExtrinsicEvents(events, blockNum)

	cb := model.ChainBlock{
		ID:             blockNum,
		Hash:           hash,
		BlockNum:       blockNum,
		ParentHash:     block.Header.ParentHash,
		StateRoot:      block.Header.StateRoot,
		ExtrinsicsRoot: block.Header.ExtrinsicsRoot,
		SpecVersion:    spec,
		Finalized:      true,
	}

	var extrinsics []model.ChainExtrinsic
	_ = util.UnmarshalAny(&extrinsics, decodeExtrinsics)
	extrinsics = s.fillExtrinsicHash(extrinsics, block.Extrinsics)
	cb.BlockTimestamp = FindOutBlockTime(extrinsics)
	err = s.createExtrinsic(ctx, txn, &cb, extrinsics, block.Extrinsics, eventMap)
	if err != nil {
		return err
	}

	if err = s.AddEvent(txn, &cb, events); err != nil {
		return err
	}
	if validator, err = s.EmitLog(txn, blockNum, logs, true, s.ValidatorsList(hash)); err != nil {
		return err
	}

	cb.Validator = validator
	cb.CodecError = codecErr != nil
	cb.ExtrinsicsCount = len(extrinsics)
	cb.EventCount = len(events)

	if err = s.dao.CreateBlock(txn, &cb); err == nil {
		s.dao.DbCommit(txn)
		s.emitBlock(ctx, &cb)
	}
	return err
}

func (s *Service) checkoutExtrinsicEvents(e []model.ChainEvent, blockNumInt uint) map[string][]model.ChainEvent {
	eventMap := make(map[string][]model.ChainEvent)
	for _, event := range e {
		extrinsicIndex := fmt.Sprintf("%d-%d", blockNumInt, event.ExtrinsicIdx)
		eventMap[extrinsicIndex] = append(eventMap[extrinsicIndex], event)
	}
	return eventMap
}

func (s *Service) GetCurrentRuntimeSpecVersion(blockNum uint) int {
	if util.CurrentRuntimeSpecVersion != 0 {
		return util.CurrentRuntimeSpecVersion
	}
	if block := s.dao.GetNearBlock(blockNum); block != nil {
		return block.SpecVersion
	}
	return -1
}

func (s *Service) GetBlockByHashJson(ctx context.Context, hash string) *model.ChainBlockJson {
	block := s.dao.GetBlockByHash(ctx, hash)
	if block == nil {
		return nil
	}
	return s.dao.BlockAsJson(ctx, block)
}

func (s *Service) EventByIndex(ctx context.Context, index string) *model.ChainEvent {
	return s.dao.GetEventByIdx(ctx, index)
}

func (s *Service) ValidatorsList(hash string) (validatorList []string) {
	validatorsRaw, _ := rpc.ReadStorage(nil, "Session", "Validators", hash)
	for _, addr := range validatorsRaw.ToStringSlice() {
		validatorList = append(validatorList, util.TrimHex(addr))
	}
	return
}

func (s *Service) fillExtrinsicHash(extrinsicList []model.ChainExtrinsic, extrinsicRaws []string) []model.ChainExtrinsic {
	for i, e := range extrinsicList {
		if e.ExtrinsicHash == "" {
			extrinsicList[i].ExtrinsicHash = util.AddHex(util.BytesToHex(hasher.HashByCryptoName(util.HexToBytes(extrinsicRaws[i]), "Blake2_256")))
		}
	}
	return extrinsicList
}
