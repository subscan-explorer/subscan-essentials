package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/itering/subscan/util"
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/itering/substrate-api-rpc/websocket"
	"github.com/panjf2000/ants/v2"
	"golang.org/x/exp/slog"
)

// FinalizedWaitingBlockCount
// Because when receive chain_finalizedHead, get block still not finalized
// so set Waiting block count to try avoid
const (
	FinalizedWaitingBlockCount = 0
	ChainNewHead               = "chain_newHead"
	ChainFinalizedHead         = "chain_finalizedHead"
	StateStorage               = "state_storage"
	BlockTime                  = 5
)

type subscription struct {
	Topic  string `json:"topic"`
	Latest int64  `json:"latest"`
}

var (
	onceFinHead     sync.Once
	subscriptionIds = []subscription{{Topic: ChainNewHead}, {Topic: ChainFinalizedHead}, {Topic: StateStorage}}
)

type SubscribeService struct {
	*Service
	newHead    chan bool
	newFinHead chan int
	done       chan struct{}
}

func (s *Service) initSubscribeService(done chan struct{}) *SubscribeService {
	return &SubscribeService{
		Service:    s,
		newHead:    make(chan bool, 1),
		newFinHead: make(chan int, 1),
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
	case ChainFinalizedHead:
		r := j.ToNewHead()
		num := util.HexToNum(r.Number)
		_ = s.updateChainMetadata(map[string]interface{}{"finalized_blockNum": util.HexToNumStr(r.Number)})
		upgradeHealth(j.Method)
		go func() {
			s.newFinHead <- int(num)
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

	p, _ := ants.NewPoolWithFunc(5, func(i interface{}) {
		blockNum := i.(BlockFinalized)
		func(bf BlockFinalized) {
			if err := s.FillBlockData(nil, bf.BlockNum, bf.Finalized); err != nil {
				logError("ChainGetBlockHash get", err)
			} else {
				s.SetHeartBeat(fmt.Sprintf("%s:heartBeat:%s", util.NetworkNode, "substrate"))
			}
		}(blockNum)
		wg.Done()
	}, ants.WithOptions(ants.Options{PanicHandler: func(c interface{}) {}}))

	defer p.Release()
	for {
		select {
		case newHead := <-s.newFinHead:
			if newHead == 0 {
				time.Sleep(BlockTime * time.Second)
				return
			}

			lastNum, _ := s.dao.GetFillFinalizedBlockNum(ctx)
			startBlock := lastNum + 1
			if lastNum == 0 {
				startBlock = lastNum
			}
			for i := startBlock; i <= int(newHead-FinalizedWaitingBlockCount); i++ {
				wg.Add(1)
				err := p.Invoke(BlockFinalized{BlockNum: i, Finalized: true})
				if err != nil {
					logError("ChainGetBlockHash get", err)
				}
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

	slog.Debug("Sending request for block", "Number", blockNum, "Finalized", finalized)
	// Block Hash
	res, err := util.SendWsRequest(conn, rpc.ChainGetBlockHash(wsBlockHash, blockNum))
	if err != nil {
		e := fmt.Errorf("ChainGetBlockHash get error %v", err)
		slog.Error("fillblockdata fail", "error", e)
		return e
	}

	blockHash, err := res.ToString()
	if err != nil || blockHash == "" {
		e := fmt.Errorf("ChainGetBlockHash get error %v", err)
		slog.Error("fillblockdata fail", "error", e)
		return e
	}
	slog.Info("Got new block", "Number", blockNum, "Hash", blockHash)

	// block
	res, err = util.SendWsRequest(conn, rpc.ChainGetBlock(wsBlock, blockHash))
	if err != nil {
		return fmt.Errorf("websocket send error: %v", err)
	}
	rpcBlock := res.ToBlock()

	// event
	res, err = util.SendWsRequest(conn, rpc.StateGetStorage(wsEvent, util.EventStorageKey, blockHash))
	if err != nil {
		return fmt.Errorf("websocket send error: %v", err)
	}
	event, _ := res.ToString()

	// runtime
	res, err = util.SendWsRequest(conn, rpc.ChainGetRuntimeVersion(wsSpec, blockHash))
	if err != nil {
		return fmt.Errorf("websocket send error: %v", err)
	}

	var specVersion int

	if r := res.ToRuntimeVersion(); r == nil {
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

	setFinalized := func() {
		if finalized {
			_ = s.dao.SaveFillAlreadyFinalizedBlockNum(context.TODO(), blockNum)
		}
	}
	// refresh finalized info for update
	// if block != nil {
	// 	// Confirm data, only set block Finalized, refresh all block data
	// 	block.ExtrinsicsRoot = rpcBlock.Block.Header.ExtrinsicsRoot
	// 	block.Hash = blockHash
	// 	block.ParentHash = rpcBlock.Block.Header.ParentHash
	// 	block.StateRoot = rpcBlock.Block.Header.StateRoot
	// 	block.Extrinsics = util.ToString(rpcBlock.Block.Extrinsics)
	// 	block.Logs = util.ToString(rpcBlock.Block.Header.Digest.Logs)
	// 	block.Event = event
	// 	_ = s.UpdateBlockData(conn, block, finalized)
	// 	return
	// }
	// for Create
	if err = s.CreateChainBlock(conn, blockHash, &rpcBlock.Block, event, specVersion, finalized); err == nil {
		_ = s.dao.SaveFillAlreadyBlockNum(context.TODO(), blockNum)
		setFinalized()
	} else {
		slog.Error("Create chain block error %v", err)
	}
	return
}

func (s *Service) updateChainMetadata(metadata map[string]interface{}) (err error) {
	c := context.TODO()
	err = s.dao.SetMetadata(c, metadata)
	return
}
