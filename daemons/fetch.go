package daemons

import (
	"errors"
	"fmt"
	"github.com/bilibili/kratos/pkg/log"
	"github.com/gorilla/websocket"
	"subscan-end/internal/service"
	"subscan-end/libs/substrate"
	"subscan-end/utiles"
	"time"
)

var newHead = make(chan bool, 1)

func fillALLBlock(srv *service.Service) {
	for {
		c, _, err := websocket.DefaultDialer.Dial(utiles.ProviderEndPoint, nil)
		if err != nil {
			log.Error("dial websocket error", err)
			time.Sleep(6 * time.Second)
			continue
		}
		<-newHead
		metadata, err := srv.GetChainMetadata()
		if err != nil {
			continue
		}
		blockNum := utiles.StringToInt(metadata["blockNum"])
		if blockNum == 0 {
			continue
		}
		alreadyBlock, _ := srv.GetAlreadyBlockNum()
		for i := alreadyBlock + 1; i <= blockNum; i++ {
			if err = fillBlockData(c, i, srv); err != nil {
				log.Error("%v", err)
				break
			}
		}
	}
}

func fillBlockData(c *websocket.Conn, i int, srv *service.Service) (err error) {
	v := &substrate.JsonRpcResult{}
	if err = c.WriteMessage(websocket.TextMessage, substrate.ChainGetBlockHash(2001, i)); err != nil {
		return errors.New(fmt.Sprintf("websocket send error: %v", err))
	}
	_ = c.ReadJSON(v)
	blockHash, _ := v.ToString()
	log.Info("Block num %d hash %s", i, blockHash)
	blockWsId := 2002
	eventWsId := 2003
	if err = c.WriteMessage(websocket.TextMessage, substrate.ChainGetBlock(blockWsId, blockHash)); err != nil {
		return errors.New(fmt.Sprintf("websocket send error: %v", err))
	}
	// event
	if err = c.WriteMessage(websocket.TextMessage, substrate.StateGetStorageAt(eventWsId, "0xcc956bdb7605e3547539f321ac2bc95c", blockHash)); err != nil {
		return errors.New(fmt.Sprintf("websocket send error: %v", err))
	}
	var (
		block *substrate.BlockResult
		event string
	)
	for i := 0; i < 2; i++ {
		v = &substrate.JsonRpcResult{}
		if err = c.ReadJSON(v); err != nil {
			return errors.New(fmt.Sprintf("websocket read json error %v", err))
		}
		if v.Id == blockWsId {
			block = v.ToBlock()
		} else {
			event, _ = v.ToString()
		}
	}
	if block == nil {
		return errors.New("nil block data")
	}
	go func() {
		_ = srv.SetAlreadyBlockNum(i)
		srv.CreateChainBlock(blockHash, &block.Block, event)
	}()
	return
}
