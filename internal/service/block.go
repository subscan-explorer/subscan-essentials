package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/pkg/log"
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/util"
	"github.com/itering/subscan/util/address"
	"github.com/itering/substrate-api-rpc"
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/itering/substrate-api-rpc/storage"
	"github.com/itering/substrate-api-rpc/websocket"
)

func (s *Service) CreateChainBlock(conn websocket.WsConn, hash string, block *rpc.Block, event string, spec int, finalized bool) (err error) {
	var (
		decodeExtrinsics []map[string]interface{}
		decodeEvent      interface{}
		logs             []storage.DecoderLog
		validator        string
	)
	c := context.TODO()

	blockNum := util.StringToInt(util.HexToNumStr(block.Header.Number))

	metadataInstant := s.getMetadataInstant(spec, hash)

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

	txn := s.dao.DbBegin()
	defer s.dao.DbRollback(txn)

	var e []model.ChainEvent
	util.UnmarshalAny(&e, decodeEvent)

	eventMap := s.checkoutExtrinsicEvents(e, blockNum)

	cb := model.ChainBlock{
		Hash:           hash,
		BlockNum:       blockNum,
		ParentHash:     block.Header.ParentHash,
		StateRoot:      block.Header.StateRoot,
		ExtrinsicsRoot: block.Header.ExtrinsicsRoot,
		Logs:           util.ToString(block.Header.Digest.Logs),
		Extrinsics:     util.ToString(block.Extrinsics),
		Event:          event,
		SpecVersion:    spec,
		Finalized:      finalized,
	}

	extrinsicsCount, blockTimestamp, extrinsicHash, extrinsicFee, err := s.createExtrinsic(c, txn, &cb, block.Extrinsics, decodeExtrinsics, eventMap)
	if err != nil {
		return err
	}
	cb.BlockTimestamp = blockTimestamp
	eventCount, err := s.AddEvent(txn, &cb, e, extrinsicHash, extrinsicFee)
	if err != nil {
		return err
	}
	if validator, err = s.EmitLog(txn, hash, blockNum, logs, finalized, s.ValidatorsList(conn, hash)); err != nil {
		return err
	}

	cb.Validator = validator
	cb.CodecError = validator == "" && blockNum != 0
	cb.ExtrinsicsCount = extrinsicsCount
	cb.EventCount = eventCount

	if err = s.dao.CreateBlock(txn, &cb); err == nil {
		s.dao.DbCommit(txn)
	}
	return err
}

func (s *Service) UpdateBlockData(conn websocket.WsConn, block *model.ChainBlock, finalized bool) (err error) {
	c := context.TODO()

	var (
		decodeEvent      interface{}
		encodeExtrinsics []string
		decodeExtrinsics []map[string]interface{}
	)

	_ = json.Unmarshal([]byte(block.Extrinsics), &encodeExtrinsics)

	spec := block.SpecVersion

	metadataInstant := s.getMetadataInstant(spec, block.Hash)

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
	util.UnmarshalAny(&e, decodeEvent)
	eventMap := s.checkoutExtrinsicEvents(e, block.BlockNum)

	txn := s.dao.DbBegin()
	defer s.dao.DbRollback(txn)

	extrinsicsCount, blockTimestamp, extrinsicHash, extrinsicFee, err := s.createExtrinsic(c, txn, block, encodeExtrinsics, decodeExtrinsics, eventMap)
	if err != nil {
		return err
	}
	block.BlockTimestamp = blockTimestamp

	eventCount, err := s.AddEvent(txn, block, e, extrinsicHash, extrinsicFee)
	if err != nil {
		return err
	}

	validator, err := s.EmitLog(txn, block.Hash, block.BlockNum, logs, finalized, s.ValidatorsList(conn, block.Hash))
	if err != nil {
		return err
	}

	if err = s.dao.UpdateEventAndExtrinsic(txn, block, eventCount, extrinsicsCount, blockTimestamp, validator, validator == "" && block.BlockNum != 0, finalized); err != nil {
		return
	}

	s.dao.DbCommit(txn)
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
	if util.CurrentRuntimeSpecVersion != 0 {
		return util.CurrentRuntimeSpecVersion
	}
	if block := s.dao.GetNearBlock(blockNum); block != nil {
		return block.SpecVersion
	}
	return -1
}

func (s *Service) GetExtrinsicList(page, row int, order string, query ...string) ([]*model.ChainExtrinsicJson, int) {
	c := context.TODO()
	list, count := s.dao.GetExtrinsicList(c, page, row, order, query...)
	var ejs []*model.ChainExtrinsicJson
	for _, extrinsic := range list {
		ejs = append(ejs, s.dao.ExtrinsicsAsJson(&extrinsic))
	}
	return ejs, count
}

func (s *Service) GetBlocksSampleByNums(page, row int) []model.SampleBlockJson {
	var blockJson []model.SampleBlockJson
	blocks := s.dao.GetBlockList(page, row)
	for _, block := range blocks {
		bj := s.BlockAsSampleJson(&block)
		blockJson = append(blockJson, *bj)
	}
	return blockJson
}

func (s *Service) GetExtrinsicByIndex(index string) *model.ExtrinsicDetail {
	c := context.TODO()
	return s.dao.GetExtrinsicsDetailByIndex(c, index)
}

func (s *Service) GetExtrinsicDetailByHash(hash string) *model.ExtrinsicDetail {
	c := context.TODO()
	return s.dao.GetExtrinsicsDetailByHash(c, hash)
}

func (s *Service) GetExtrinsicByHash(hash string) *model.ChainExtrinsic {
	c := context.TODO()
	return s.dao.GetExtrinsicsByHash(c, hash)
}

func (s *Service) GetBlockByHashJson(hash string) *model.ChainBlockJson {
	c := context.TODO()
	block := s.dao.GetBlockByHash(c, hash)
	if block == nil {
		return nil
	}
	return s.dao.BlockAsJson(c, block)
}

func (s *Service) EventByIndex(index string) *model.ChainEvent {
	return s.dao.GetEventByIdx(index)
}

func (s *Service) GetBlockByNum(num int) *model.ChainBlockJson {
	c := context.TODO()
	block := s.dao.GetBlockByNum(num)
	if block == nil {
		return nil
	}
	return s.dao.BlockAsJson(c, block)
}

func (s *Service) GetBlockByHash(hash string) *model.ChainBlock {
	c := context.TODO()
	block := s.dao.GetBlockByHash(c, hash)
	if block == nil {
		return nil
	}
	return block
}

func (s *Service) BlockAsSampleJson(block *model.ChainBlock) *model.SampleBlockJson {
	b := model.SampleBlockJson{
		BlockNum:        block.BlockNum,
		BlockTimestamp:  block.BlockTimestamp,
		Hash:            block.Hash,
		EventCount:      block.EventCount,
		ExtrinsicsCount: block.ExtrinsicsCount,
		Validator:       address.SS58Address(block.Validator),
		Finalized:       block.Finalized,
	}
	return &b
}

func (s *Service) GetCurrentBlockNum(c context.Context) (uint64, error) {
	return s.dao.GetBestBlockNum(c)
}

func (s *Service) ValidatorsList(conn websocket.WsConn, hash string) (validatorList []string) {
	validatorsRaw, _ := rpc.ReadStorage(conn, "Session", "Validators", hash)
	for _, addr := range validatorsRaw.ToStringSlice() {
		validatorList = append(validatorList, util.TrimHex(addr))
	}
	return
}
