package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/pkg/log"
	"github.com/itering/subscan/util"
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/itering/substrate-api-rpc/websocket"
	"github.com/panjf2000/ants/v2"
	"sync"
	"time"
)

// FinalizedWaitingBlockCount
// Because when receive chain_finalizedHead, get block still not finalized
// so set Waiting block count to try avoid
const (
	FinalizedWaitingBlockCount = 3
	ChainNewHead               = "chain_newHead"
	ChainFinalizedHead         = "chain_finalizedHead"
	StateStorage               = "state_storage"
	BlockTime                  = 6
)

type subscription struct {
	Topic  string `json:"topic"`
	Latest int64  `json:"latest"`
}

var (
	onceNewHead, onceFinHead sync.Once
	subscriptionIds          = []subscription{{Topic: ChainNewHead}, {Topic: ChainFinalizedHead}, {Topic: StateStorage}}
)

type SubscribeService struct {
	*Service
	newHead    chan bool
	newFinHead chan bool
	done       chan struct{}
}

func (s *Service) initSubscribeService(done chan struct{}) *SubscribeService {
	return &SubscribeService{
		Service:    s,
		newHead:    make(chan bool, 1),
		newFinHead: make(chan bool, 1),
		done:       done,
	}
}

func (s *SubscribeService) parser(message []byte) (err error) {
	upgradeHealth := func(topic string) {
		for index, subscript := range subscriptionIds {
			if subscript.Topic == topic {
				subscriptionIds[index].Latest = time.Now().Unix()
			}
		}
	}

	var j rpc.JsonRpcResult
	if err = json.Unmarshal(message, &j); err != nil {
		return err
	}

	switch j.Id {
	case runtimeVersion:
		r := j.ToRuntimeVersion()
		_ = s.regRuntimeVersion(r.ImplName, r.SpecVersion)
		_ = s.updateChainMetadata(map[string]interface{}{"implName": r.ImplName, "specVersion": r.SpecVersion})
		util.CurrentRuntimeSpecVersion = r.SpecVersion
		return
	}

	switch j.Method {
	case ChainNewHead:
		r := j.ToNewHead()
		_ = s.updateChainMetadata(map[string]interface{}{"blockNum": util.HexToNumStr(r.Number)})
		upgradeHealth(j.Method)
		go func() {
			s.newHead <- true
			onceNewHead.Do(func() {
				go s.subscribeFetchBlock()
			})
		}()
	case ChainFinalizedHead:
		r := j.ToNewHead()
		_ = s.updateChainMetadata(map[string]interface{}{"finalized_blockNum": util.HexToNumStr(r.Number)})
		upgradeHealth(j.Method)
		go func() {
			s.newFinHead <- true
			onceFinHead.Do(func() {
				go s.subscribeFetchBlock()
			})
		}()
	case StateStorage:
		upgradeHealth(j.Method)
	default:
		return
	}
	return
}

func (s *SubscribeService) subscribeFetchBlock() {
	var wg sync.WaitGroup
	ctx := context.TODO()
	type BlockFinalized struct {
		BlockNum  int  `json:"block_num"`
		Finalized bool `json:"finalized"`
	}

	p, _ := ants.NewPoolWithFunc(10, func(i interface{}) {
		blockNum := i.(BlockFinalized)
		func(bf BlockFinalized) {
			if err := s.FillBlockData(nil, bf.BlockNum, bf.Finalized); err != nil {
				log.Error("ChainGetBlockHash get error %v", err)
			} else {
				s.SetHeartBeat(fmt.Sprintf("%s:heartBeat:%s", util.NetworkNode, "substrate"))
			}
		}(blockNum)
		wg.Done()
	}, ants.WithOptions(ants.Options{PanicHandler: func(c interface{}) {}}))

	defer p.Release()
	for {
		select {
		case <-s.newHead:
			best, err := s.dao.GetBestBlockNum(ctx)
			if err != nil || best == 0 {
				time.Sleep(BlockTime * time.Second)
				return
			}
			lastNum, _ := s.dao.GetFillBestBlockNum(ctx)
			finalizedBlock, _ := s.dao.GetFinalizedBlockNum(ctx)

			startBlock := lastNum + 1
			if lastNum == 0 {
				startBlock = lastNum
			}

			for i := startBlock; i <= int(best); i++ {
				wg.Add(1)
				_ = p.Invoke(BlockFinalized{
					BlockNum:  i,
					Finalized: finalizedBlock >= FinalizedWaitingBlockCount && uint64(i) <= finalizedBlock-FinalizedWaitingBlockCount,
				})
			}
			wg.Wait()
		case <-s.newFinHead:
			final, err := s.dao.GetFinalizedBlockNum(context.TODO())
			if err != nil || final == 0 {
				time.Sleep(BlockTime * time.Second)
				return
			}

			lastNum, _ := s.dao.GetFillFinalizedBlockNum(ctx)
			startBlock := lastNum + 1
			if lastNum == 0 {
				startBlock = lastNum
			}

			for i := startBlock; i <= int(final-FinalizedWaitingBlockCount); i++ {
				wg.Add(1)
				_ = p.Invoke(BlockFinalized{BlockNum: i, Finalized: true})
			}
			wg.Wait()
		case <-s.done:
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

func (s *Service) FillBlockData(conn websocket.WsConn, blockNum int, finalized bool) (err error) {
	block := s.dao.GetBlockByNum(blockNum)
	if block != nil && block.Finalized && !block.CodecError {
		return nil
	}

	v := &rpc.JsonRpcResult{}

	// Block Hash
	if err = websocket.SendWsRequest(conn, v, rpc.ChainGetBlockHash(wsBlockHash, blockNum)); err != nil {
		return fmt.Errorf("websocket send error: %v", err)
	}
	blockHash, err := v.ToString()
	if err != nil || blockHash == "" {
		return fmt.Errorf("ChainGetBlockHash get error %v", err)
	}
	log.Info("Block num %d hash %s", blockNum, blockHash)

	// block
	if err = websocket.SendWsRequest(conn, v, rpc.ChainGetBlock(wsBlock, blockHash)); err != nil {
		return fmt.Errorf("websocket send error: %v", err)
	}
	rpcBlock := v.ToBlock()

	// event
	if err = websocket.SendWsRequest(conn, v, rpc.StateGetStorage(wsEvent, util.EventStorageKey, blockHash)); err != nil {
		return fmt.Errorf("websocket send error: %v", err)
	}
	event, _ := v.ToString()

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

	if rpcBlock == nil || specVersion == -1 {
		return errors.New("nil block data")
	}

	var setFinalized = func() {
		if finalized {
			_ = s.dao.SaveFillAlreadyFinalizedBlockNum(context.TODO(), blockNum)
		}
	}
	// refresh finalized info for update
	if block != nil {
		// Confirm data, only set block Finalized
		if block.Hash == blockHash && block.ExtrinsicsRoot == rpcBlock.Block.Header.ExtrinsicsRoot && block.Event == event && !block.CodecError && finalized {
			s.dao.SetBlockFinalized(block)
		} else {
			// refresh all block data
			block.ExtrinsicsRoot = rpcBlock.Block.Header.ExtrinsicsRoot
			block.Hash = blockHash
			block.ParentHash = rpcBlock.Block.Header.ParentHash
			block.StateRoot = rpcBlock.Block.Header.StateRoot

			block.Extrinsics = util.ToString(rpcBlock.Block.Extrinsics)
			block.Logs = util.ToString(rpcBlock.Block.Header.Digest.Logs)
			block.Event = event

			_ = s.UpdateBlockData(conn, block, finalized)
		}
		setFinalized()
		return
	}
	// for Create
	if err = s.CreateChainBlock(conn, blockHash, &rpcBlock.Block, event, specVersion, finalized); err == nil {
		_ = s.dao.SaveFillAlreadyBlockNum(context.TODO(), blockNum)
		setFinalized()
	} else {
		log.Error("Create chain block error %v", err)
	}
	return
}

func (s *Service) updateChainMetadata(metadata map[string]interface{}) (err error) {
	c := context.TODO()
	err = s.dao.SetMetadata(c, metadata)
	return
}
