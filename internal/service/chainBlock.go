package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bilibili/kratos/pkg/log"
	"subscan-end/internal/model"
	"subscan-end/libs/substrate"
	"subscan-end/libs/substrate/protos/codec_protos"
	"subscan-end/utiles"
	"subscan-end/utiles/pusher/dingding"
)

func (s *Service) GetAlreadyBlockNum() (int, error) {
	c := context.TODO()
	return s.dao.GetFillAlreadyBlockNum(c)
}

func (s *Service) SetAlreadyBlockNum(num int) error {
	c := context.TODO()
	return s.dao.SaveFillAlreadyBlockNum(c, num)
}

func (s *Service) GetRepairBlockBlockNum() (int, error) {
	c := context.TODO()
	return s.dao.GetRepairBlockBlockNum(c)
}

func (s *Service) SetRepairBlockBlockNum(num int) error {
	c := context.TODO()
	return s.dao.SaveRepairBlockBlockNum(c, num)
}

func (s *Service) GetBlockNumArr(start, end int) []int {
	c := context.TODO()
	return s.dao.GetBlockNumArr(c, start, end)
}

func (s *Service) GetBlockList(page, row int) []model.ChainBlock {
	c := context.TODO()
	return s.dao.GetBlockList(c, page, row)
}

func (s *Service) GetBlockByHashJson(hash string) *model.ChainBlockJson {
	c := context.TODO()
	block := s.dao.BlockByHash(c, hash)
	if block == nil {
		return nil
	}
	return s.dao.BlockAsJson(c, block)
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

func (s *Service) GetBlockFixDataList() *[]model.ChainBlock {
	c := context.TODO()
	return s.dao.GetAllBlocksNeedFix(c)
}

func (s *Service) GetBlocksSampleByNums(page, row int) *[]model.SampleBlockJson {
	c := context.TODO()
	blocks := s.dao.GetBlockList(c, page, row)
	var blockJson []model.SampleBlockJson
	for _, block := range blocks {
		blockJson = append(blockJson, *s.dao.BlockAsSampleJson(c, &block))
	}
	return &blockJson
}

func (s *Service) CreateChainBlock(hash string, block *substrate.Block, event string) {
	c := context.TODO()
	extrinsicB, _ := json.Marshal(block.Extrinsics)
	metadataVersion := substrate.GetNowMetadataVersion(utiles.StringToInt(utiles.HexToNumStr(block.Header.Number)))
	var err error
	codecError := false
	decodeExtrinsics := ""
	if string(extrinsicB) != "" {
		decodeExtrinsics, err = codec_protos.DecodeExtrinsic(string(extrinsicB), metadataVersion)
		if err != nil {
			go dingding.DingClient.Push("text", "ERR: Decode Extrinsics get error", string(extrinsicB), err.Error())
			log.Error("ERR: Decode Extrinsics get error ", err)
			codecError = true
		}
	}
	decodeEvent := ""
	if event != "" {
		decodeEvent, err = codec_protos.DecodeEvent(event, metadataVersion)
		if err != nil {
			go dingding.DingClient.Push("text", "ERR: Decode Event get error", event, err.Error())
			log.Error("ERR: Decode Event get error ", err)
			codecError = true
		}
	}
	logByte, _ := json.Marshal(block.Header.Digest.Logs)
	logStr := string(logByte)
	decodeLog := ""
	if logStr != "" {
		decodeLog, err = codec_protos.DecodeLog(logStr, metadataVersion)
		if err != nil {
			go dingding.DingClient.Push("text", "ERR: Decode Log get error", logStr, err.Error())
			log.Error("ERR: Decode Log get error ", err)
			codecError = true
		}
	}
	txn := s.dao.DbBegin()
	defer txn.DbRollback()
	var (
		eventCount      int
		extrinsicsCount int
		blockTimestamp  int
		validator       string
	)
	if codecError == false {
		var hashMap map[string]string
		var e []model.ChainEvent
		_ = json.Unmarshal([]byte(decodeEvent), &e)
		blockNumInt := utiles.StringToInt(utiles.HexToNumStr(block.Header.Number))
		successMap := s.checkoutExtrinsicSuccess(&e, blockNumInt)
		extrinsicsCount, blockTimestamp, hashMap = s.createExtrinsic(c, txn, block.Header.Number, decodeExtrinsics, successMap)
		eventCount = s.AddEvent(c, txn, blockNumInt, blockTimestamp, hash, e, hashMap)
		validator = s.EmitLog(c, txn, hash, block.Header.Number, decodeLog)
	}
	_ = s.dao.CreateBlock(c, txn, hash, block, event, decodeEvent, decodeExtrinsics, logStr, decodeLog, validator, eventCount, extrinsicsCount, blockTimestamp, codecError)
	txn.DbCommit()
}

func (s *Service) RepairBlockData(block *model.ChainBlock) {
	c := context.TODO()
	metadataVersion := substrate.GetNowMetadataVersion(block.BlockNum)
	var (
		decodeEvent string
		err         error
	)
	if block.Event != "" {
		decodeEvent, err = codec_protos.DecodeEvent(block.Event, metadataVersion)
		if err != nil {
			log.Error("ERR: Decode Event get error ", err)
			return
		}
	}
	decodeExtrinsics, err := codec_protos.DecodeExtrinsic(block.Extrinsics, metadataVersion)
	if err != nil {
		log.Error("ERR: Decode Extrinsics get error ", err)
		return
	}
	decodeLog, err := codec_protos.DecodeLog(block.Logs, metadataVersion)
	if err != nil {
		log.Error("ERR: Decode Logs get error ", err)
		return
	}
	txn := s.dao.DbBegin()
	defer txn.DbRollback()
	var e []model.ChainEvent
	_ = json.Unmarshal([]byte(decodeEvent), &e)
	successMap := s.checkoutExtrinsicSuccess(&e, block.BlockNum)
	extrinsicsCount, blockTimestamp, hashMap := s.createExtrinsic(c, txn, utiles.IntToHex(block.BlockNum), decodeExtrinsics, successMap)
	eventCount := s.AddEvent(c, txn, block.BlockNum, blockTimestamp, block.Hash, e, hashMap)
	validator := s.EmitLog(c, txn, block.Hash, utiles.IntToHex(block.BlockNum), decodeLog)
	if err = s.dao.UpdateEventAndExtrinsic(c, txn, block, decodeEvent, decodeExtrinsics, decodeLog, eventCount, extrinsicsCount, blockTimestamp, validator); err != nil {
		return
	}
	txn.DbCommit()
}

func (s *Service) checkoutExtrinsicSuccess(e *[]model.ChainEvent, blockNumInt int) map[string]bool {
	successMap := make(map[string]bool)
	for _, event := range *e {
		if event.ModuleId == "system" {
			eventIndex := fmt.Sprintf("%d-%d", blockNumInt, event.ExtrinsicIdx)
			successMap[eventIndex] = event.EventId != "ExtrinsicFailed"
		}
	}
	return successMap
}
