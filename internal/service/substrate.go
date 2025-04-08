package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/itering/substrate-api-rpc/model"
	"sync"
	"time"

	"log"

	"github.com/itering/subscan/util"
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/itering/substrate-api-rpc/websocket"
)

const (
	FinalizedWaitingBlockCount = 2
	BlockTime                  = 6
	ChainFinalizedHead         = "chain_finalizedHead"
	StateRuntimeVersion        = "state_runtimeVersion"
)

var (
	onceFinHead sync.Once
)

type SubscribeService struct {
	*Service
	newFinHead        chan bool
	lastBlock         int64
	finalizedBlockNum int64
}

func (s *Service) initSubscribeService() *SubscribeService {
	return &SubscribeService{
		Service:    s,
		newFinHead: make(chan bool, 1),
	}
}

func (s *SubscribeService) parser(message []byte) (err error) {
	var j model.JsonRpcResult
	if err = json.Unmarshal(message, &j); err != nil {
		return err
	}
	switch j.Method {
	case ChainFinalizedHead:
		r := j.ToNewHead()
		_ = s.updateChainMetadata(map[string]interface{}{"finalized_blockNum": util.HexToNumStr(r.Number)})
		s.finalizedBlockNum = util.U256(r.Number).Int64()
		s.newFinHead <- true
	case StateRuntimeVersion:
		r := j.ToRuntimeVersion()
		_ = s.regRuntimeVersion(r.ImplName, r.SpecVersion)
		_ = s.updateChainMetadata(map[string]interface{}{"implName": r.ImplName, "specVersion": r.SpecVersion})
		util.CurrentRuntimeSpecVersion = r.SpecVersion
	default:
		return
	}
	return
}

func (s *SubscribeService) subscribeFetchBlock(ctx context.Context) {
	for {
		select {
		case <-s.newFinHead:

			if s.finalizedBlockNum == 0 {
				time.Sleep(BlockTime * time.Second)
				return
			}

			lastNum, _ := s.dao.GetFillFinalizedBlockNum(ctx)
			startBlock := int64(lastNum)
			if s.lastBlock > 0 {
				startBlock = s.lastBlock + 1
			}
			for i := startBlock; i <= s.finalizedBlockNum-FinalizedWaitingBlockCount; i++ {
				// mq.Instant.Publish("block", "block", map[string]interface{}{"block_num": i})
				util.Logger().Info(fmt.Sprintf("Publish block num %d", i))
				s.lastBlock = i
			}
		case <-ctx.Done():
			return
		}
	}
}

const (
	wsBlockHash = iota + 1
	wsBlock
	wsEvent
	wsSpec
)

func (s *Service) FillBlockData(ctx context.Context, blockNum uint) (err error) {
	block := s.dao.GetBlockByNum(ctx, blockNum)
	if block != nil && block.Finalized && !block.CodecError {
		return nil
	}

	conn := s.dbStorage.RPCPool().Conn

	v := &model.JsonRpcResult{}

	// Block Hash
	if err = websocket.SendWsRequest(conn, v, rpc.ChainGetBlockHash(wsBlockHash, int(blockNum))); err != nil {
		return fmt.Errorf("websocket send error: %v", err)
	}
	blockHash, err := v.ToString()
	if err != nil || blockHash == "" {
		return fmt.Errorf("ChainGetBlockHash get error %v", err)
	}
	util.Logger().Info(fmt.Sprintf("Block num %d hash %s", blockNum, blockHash))

	// block
	if err = websocket.SendWsRequest(conn, v, rpc.ChainGetBlock(wsBlock, blockHash)); err != nil {
		return fmt.Errorf("websocket send error: %v", err)
	}
	rpcBlock := v.ToBlock()
	if rpcBlock == nil {
		return errors.New("nil block data")
	}

	// event
	if err = websocket.SendWsRequest(conn, v, rpc.StateGetStorage(wsEvent, util.EventStorageKey, blockHash)); err != nil {
		return fmt.Errorf("websocket send error: %v", err)
	}
	event, _ := v.ToString()
	if event == "" {
		return errors.New("nil event data")
	}

	// runtime
	if err = websocket.SendWsRequest(conn, v, rpc.ChainGetRuntimeVersion(wsSpec, blockHash)); err != nil {
		return fmt.Errorf("websocket send error: %v", err)
	}

	var specVersion int

	if r := v.ToRuntimeVersion(); r == nil {
		specVersion = s.GetCurrentRuntimeSpecVersion(blockNum)
	} else {
		specVersion = r.SpecVersion
		_ = s.regRuntimeVersion(r.ImplName, specVersion, blockHash)
	}

	if specVersion > util.CurrentRuntimeSpecVersion {
		util.CurrentRuntimeSpecVersion = specVersion
	}

	if specVersion == -1 {
		return errors.New("nil runtime version")
	}

	var setFinalized = func() {
		_ = s.dao.SaveFillAlreadyFinalizedBlockNum(context.TODO(), int(blockNum))
	}
	// for Create
	if err = s.CreateChainBlock(ctx, conn, blockHash, &rpcBlock.Block, event, specVersion); err == nil {
		_ = s.dao.SaveFillAlreadyBlockNum(ctx, int(blockNum))
		setFinalized()
	} else {
		log.Printf("Create chain block error %v", err)
	}
	return
}

func (s *Service) updateChainMetadata(metadata map[string]interface{}) (err error) {
	c := context.TODO()
	err = s.dao.SetMetadata(c, metadata)
	return
}
