package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-kratos/kratos/pkg/log"
	"github.com/itering/subscan/internal/model"
	"github.com/itering/subscan/internal/substrate"
	"github.com/itering/subscan/internal/substrate/metadata"
	"github.com/itering/subscan/internal/substrate/rpc"
	"github.com/itering/subscan/internal/substrate/storage"
	"github.com/itering/subscan/internal/util"
	"github.com/shopspring/decimal"
)

func (s *Service) GetAlreadyBlockNum() (int, error) {
	c := context.TODO()
	return s.dao.GetFillAlreadyBlockNum(c)
}

func (s *Service) GetFillFinalizedBlockNum() (int, error) {
	c := context.TODO()
	return s.dao.GetFillFinalizedBlockNum(c)
}

func (s *Service) SetAlreadyBlockNum(num int) error {
	c := context.TODO()
	return s.dao.SaveFillAlreadyBlockNum(c, num)
}

func (s *Service) GetBlockByHash(hash string) *model.ChainBlock {
	c := context.TODO()
	block := s.dao.BlockByHash(c, hash)
	if block == nil {
		return nil
	}
	return block
}

func (s *Service) GetBlockByNum(num int) *model.ChainBlockJson {
	c := context.TODO()
	block := s.dao.Block(c, num)
	if block == nil {
		return nil
	}
	return s.dao.BlockAsJson(c, block)
}

func (s *Service) GetBlockFixDataList(option ...string) []int {
	return s.dao.GetAllBlocksNeedFix(option...)
}

func (s *Service) CreateChainBlock(hash string, block *rpc.Block, event string, spec int, finalized bool) (err error) {
	var decodeExtrinsics []interface{}
	var decodeEvent interface{}
	var logs []storage.DecoderLog
	c := context.TODO()

	blockNum := util.StringToInt(util.HexToNumStr(block.Header.Number))

	metadataInstant := s.getMetadataInstant(spec)
	// Extrinsic
	decodeExtrinsics, err = substrate.DecodeExtrinsic(block.Extrinsics, metadataInstant, spec)
	if err != nil {
		log.Error("%v", err)
	}

	// event
	if err == nil {
		decodeEvent, err = substrate.DecodeEvent(event, metadataInstant, spec)
		if err != nil {
			log.Error("%v", err)
		}
	}

	// log
	if err == nil {
		logs, err = substrate.DecodeLogDigest(block.Header.Digest.Logs)
		if err != nil {
			log.Error("%v", err)
		}
	}

	var (
		eventCount, extrinsicsCount, blockTimestamp int
		validator                                   string
		extrinsicHash                               map[string]string
		extrinsicFee                                map[string]decimal.Decimal
	)

	validatorList, _ := rpc.GetValidatorFromSub(nil, hash)
	txn := s.dao.DbBegin()
	defer txn.DbRollback()

	var e []model.ChainEvent
	util.UnmarshalToAnything(&e, decodeEvent)

	eventMap := s.checkoutExtrinsicEvents(e, blockNum)

	if extrinsicsCount, blockTimestamp, extrinsicHash, extrinsicFee, err = s.createExtrinsic(c, txn, blockNum, block.Extrinsics, decodeExtrinsics, eventMap, finalized, spec); err != nil {
		return err
	}

	if eventCount, err = s.AddEvent(c, txn, blockNum, blockTimestamp, hash, e, extrinsicHash, finalized, spec, extrinsicFee); err != nil {
		return err
	}
	if validator, err = s.EmitLog(c, txn, blockNum, logs, validatorList, finalized); err != nil {
		return err
	}

	codecError := validator == ""

	if err = s.dao.CreateBlock(c, txn, hash, block, event, util.InterfaceToString(block.Header.Digest.Logs), validator, eventCount, extrinsicsCount, blockTimestamp, codecError, spec, finalized); err == nil {
		txn.DbCommit()
	}
	return err
}

func (s *Service) UpdateBlockData(block *model.ChainBlock, finalized bool) (err error) {
	c := context.TODO()

	var (
		decodeEvent      interface{}
		encodeExtrinsics []string
		decodeExtrinsics []interface{}
	)

	_ = json.Unmarshal([]byte(block.Extrinsics), &encodeExtrinsics)

	spec := block.SpecVersion

	metadataInstant := s.getMetadataInstant(spec)

	// Event
	decodeEvent, err = substrate.DecodeEvent(block.Event, metadataInstant, spec)
	if err != nil {
		fmt.Println("ERR: Decode Event get error ", err)
		return
	}

	// Extrinsic
	decodeExtrinsics, err = substrate.DecodeExtrinsic(encodeExtrinsics, metadataInstant, spec)
	if err != nil {
		fmt.Println("ERR: Decode Extrinsic get error ", err)
		return
	}

	// Log
	var rawList []string
	_ = json.Unmarshal([]byte(block.Logs), &rawList)
	logs, err := substrate.DecodeLogDigest(rawList)
	if err != nil {
		fmt.Println("ERR: Decode Logs get error ", err)
		return
	}

	var e []model.ChainEvent
	util.UnmarshalToAnything(&e, decodeEvent)

	eventMap := s.checkoutExtrinsicEvents(e, block.BlockNum)

	validatorList, _ := rpc.GetValidatorFromSub(nil, block.Hash)

	txn := s.dao.DbBegin()
	defer txn.DbRollback()

	extrinsicsCount, blockTimestamp, extrinsicHash, extrinsicFee, err := s.createExtrinsic(c, txn, block.BlockNum, encodeExtrinsics, decodeExtrinsics, eventMap, finalized, spec)
	if err != nil {
		return err
	}

	eventCount, err := s.AddEvent(c, txn, block.BlockNum, blockTimestamp, block.Hash, e, extrinsicHash, finalized, spec, extrinsicFee)
	if err != nil {
		return err
	}

	validator, err := s.EmitLog(c, txn, block.BlockNum, logs, validatorList, finalized)
	if err != nil {
		return err
	}

	if err = s.dao.UpdateEventAndExtrinsic(c, txn, block, eventCount, extrinsicsCount, blockTimestamp, validator, validator == "", finalized); err != nil {
		return
	}

	txn.DbCommit()
	return
}

func (s *Service) checkoutExtrinsicEvents(e []model.ChainEvent, blockNumInt int) map[string][]model.ChainEvent {
	eventMap := make(map[string][]model.ChainEvent)
	for _, event := range e {
		extrinsicIndex := fmt.Sprintf("%d-%d", blockNumInt, event.ExtrinsicIdx)
		eventMap[extrinsicIndex] = append(eventMap[extrinsicIndex], event)
	}
	return eventMap
}

func (s *Service) GetCurrentRuntimeSpecVersion(blockNum int) int {
	c := context.TODO()
	if substrate.CurrentRuntimeSpecVersion != 0 {
		return substrate.CurrentRuntimeSpecVersion
	}
	if block := s.dao.GetNearBlock(c, blockNum); block != nil {
		return block.SpecVersion
	}
	return 0
}

func (s *Service) getMetadataInstant(spec int) *metadata.MetadataType {
	metadataInstant, ok := metadata.RuntimeMetadata[spec]
	if !ok {
		metadataInstant = metadata.Init(s.dao.RuntimeVersionRaws(spec), spec)
	}
	return metadataInstant
}
