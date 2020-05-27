package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/freehere107/go-workers"
	"github.com/go-kratos/kratos/pkg/log"
	"github.com/itering/subscan/internal/substrate"
	"github.com/itering/subscan/internal/substrate/rpc"
	"github.com/itering/subscan/internal/substrate/websocket"
	"github.com/itering/subscan/internal/util"
	"github.com/panjf2000/ants/v2"
	"sync"
	"time"
)

// FinalizedWaitingBlockCount
// Because when receive chain_finalizedHead, get block still not finalized
// so set Waiting block count to try avoid
const FinalizedWaitingBlockCount = 5

var (
	onceNewHead, onceFinHead sync.Once
)

type SubscribeService struct {
	*Service
	newHead    chan bool
	newFinHead chan bool
}

func (s *Service) InitSubscribeService() *SubscribeService {
	return &SubscribeService{
		Service:    s,
		newHead:    make(chan bool, 1),
		newFinHead: make(chan bool, 1),
	}
}

func (s *SubscribeService) ParserSubscribe(message []byte) {
	var j rpc.JsonRpcResult
	if err := json.Unmarshal(message, &j); err != nil {
		return
	}
	if j.Id == 1 { // runtime version
		r := j.ToRuntimeVersion()
		_ = s.regRuntimeVersion(r.ImplName, r.SpecVersion)
		_ = s.UpdateChainMetadata(map[string]interface{}{"implName": r.ImplName, "specVersion": r.SpecVersion})
		substrate.CurrentRuntimeSpecVersion = r.SpecVersion
	}
	switch j.Method {
	case substrate.ChainNewHead:
		r := j.ToNewHead()
		_ = s.UpdateChainMetadata(map[string]interface{}{"blockNum": util.HexToNumStr(r.Number)})
		go func() {
			s.newHead <- true
			onceNewHead.Do(func() {
				go s.subscribeFetchBlock()
			})
		}()
	case substrate.ChainFinalizedHead:
		r := j.ToNewHead()
		_ = s.UpdateChainMetadata(map[string]interface{}{"finalized_blockNum": util.HexToNumStr(r.Number)})
		go func() {
			s.newFinHead <- true
			onceFinHead.Do(func() {
				go s.subscribeFetchBlock()
			})
		}()
	default:
		return
	}
}

func (s *SubscribeService) subscribeFetchBlock() {
	var wg sync.WaitGroup

	type BlockFinalized struct {
		BlockNum  int  `json:"block_num"`
		Finalized bool `json:"finalized"`
	}
	p, _ := ants.NewPoolWithFunc(10, func(i interface{}) {
		blockNum := i.(BlockFinalized)
		func(bf BlockFinalized) {
			if err := s.FillBlockData(bf.BlockNum, bf.Finalized); err != nil {
				_, _ = workers.EnqueueWithOptions("block", "block",
					map[string]interface{}{"block_num": bf.BlockNum, "finalized": bf.Finalized},
					workers.EnqueueOptions{RetryCount: 2})
			} else {
				s.SetHeartBeat(fmt.Sprintf("%s:heartBeat:%s", util.NetworkNode, "substrate"))
			}
		}(blockNum)
		wg.Done()
	})
	defer p.Release()
	for {
		select {
		case <-s.newHead:
			blockNum, err := s.GetCurrentBlockNum(context.TODO())
			if err != nil || blockNum == 0 {
				time.Sleep(substrate.BlockTime * time.Second)
				return
			}
			alreadyBlock, _ := s.GetAlreadyBlockNum()
			finalizedBlock, _ := s.GetFinalizedBlockNum(context.TODO())

			startBlock := alreadyBlock + 1
			if alreadyBlock == 0 {
				startBlock = 0
			}
			for i := startBlock; i <= int(blockNum); i++ {
				wg.Add(1)
				_ = p.Invoke(BlockFinalized{BlockNum: i, Finalized: finalizedBlock >= FinalizedWaitingBlockCount && uint64(i) <= finalizedBlock-FinalizedWaitingBlockCount})
			}
			wg.Wait()
		case <-s.newFinHead:
			blockNum, err := s.GetFinalizedBlockNum(context.TODO())
			if err != nil || blockNum == 0 {
				time.Sleep(substrate.BlockTime * time.Second)
				return
			}
			alreadyBlock, _ := s.GetFillFinalizedBlockNum()
			startBlock := alreadyBlock + 1
			if alreadyBlock == 0 {
				startBlock = 0
			}
			for i := startBlock; i <= int(blockNum-FinalizedWaitingBlockCount); i++ {
				wg.Add(1)
				_ = p.Invoke(BlockFinalized{BlockNum: i, Finalized: true})
			}
			wg.Wait()
		}
	}
}

const (
	wsBlockHash = iota + 1
	wsBlock
	wsEvent
	wsSpec
)

func (s *Service) FillBlockData(blockNum int, finalized bool) (err error) {
	block := s.dao.GetBlockByNum(context.TODO(), blockNum)
	if block != nil && block.Finalized && !block.CodecError {
		return nil
	}
	const websocketTextMessage = 1

	wsPool, err := websocket.Init()
	if err != nil {
		return
	}
	c := wsPool.Conn
	defer c.Close()

	v := &rpc.JsonRpcResult{}
	if err = c.WriteMessage(websocketTextMessage, rpc.ChainGetBlockHash(wsBlockHash, blockNum)); err != nil {
		return fmt.Errorf("websocket send error: %v", err)
	}
	_ = c.ReadJSON(v)
	blockHash, err := v.ToString()

	if err != nil {
		return fmt.Errorf("ChainGetBlockHash get error %v", err)
	}

	log.Info("Block num %d hash %s", blockNum, blockHash)
	if err = c.WriteMessage(websocketTextMessage, rpc.ChainGetBlock(wsBlock, blockHash)); err != nil {
		return fmt.Errorf("websocket send error: %v", err)
	}
	// event
	if err = c.WriteMessage(websocketTextMessage, rpc.StateGetStorage(wsEvent, substrate.EventStorageKey, blockHash)); err != nil {
		return fmt.Errorf("websocket send error: %v", err)
	}
	if err = c.WriteMessage(websocketTextMessage, rpc.ChainGetRuntimeVersion(wsSpec, blockHash)); err != nil {
		return fmt.Errorf("websocket send error: %v", err)
	}

	var (
		rpcBlock       *rpc.BlockResult
		event          string
		runtimeVersion = -1
	)
	for i := 0; i < 3; i++ {
		v = &rpc.JsonRpcResult{}
		if err = c.ReadJSON(v); err != nil {
			return fmt.Errorf("websocket read json error %v", err)
		}
		switch v.Id {
		case wsBlock:
			rpcBlock = v.ToBlock()
		case wsEvent:
			event, _ = v.ToString()
		default:
			if r := v.ToRuntimeVersion(); r == nil {
				runtimeVersion = s.GetCurrentRuntimeSpecVersion(i)
			} else {
				runtimeVersion = r.SpecVersion
				if runtimeVersion > substrate.CurrentRuntimeSpecVersion && substrate.CurrentRuntimeSpecVersion >= 0 {
					substrate.CurrentRuntimeSpecVersion = runtimeVersion
					_ = s.regRuntimeVersion(r.ImplName, runtimeVersion)
				}
			}
		}
	}
	if rpcBlock == nil || runtimeVersion == -1 {
		return errors.New("nil block data")
	}

	var setFinalized = func() {
		if finalized {
			_ = s.dao.SaveFillAlreadyFinalizedBlockNum(context.TODO(), blockNum)
		}
	}

	// refresh finalized info for update
	if block != nil {
		// Confirm data, set block Finalized
		if block.ExtrinsicsRoot == rpcBlock.Block.Header.ExtrinsicsRoot && block.Event == event && !block.CodecError && finalized {
			s.dao.SetBlockFinalized(block)
		} else {
			block.ExtrinsicsRoot = rpcBlock.Block.Header.ExtrinsicsRoot
			block.Hash = blockHash
			block.ParentHash = rpcBlock.Block.Header.ParentHash
			block.StateRoot = rpcBlock.Block.Header.StateRoot

			extrinsicB, _ := json.Marshal(rpcBlock.Block.Extrinsics)
			block.Extrinsics = string(extrinsicB)

			logByte, _ := json.Marshal(rpcBlock.Block.Header.Digest.Logs)
			block.Logs = string(logByte)

			block.Event = event

			_ = s.UpdateBlockData(block, finalized)
		}
		setFinalized()
		return
	}

	// for Create
	if err = s.CreateChainBlock(blockHash, &rpcBlock.Block, event, runtimeVersion, finalized); err == nil {
		_ = s.SetAlreadyBlockNum(blockNum)
		setFinalized()
	}
	return
}
