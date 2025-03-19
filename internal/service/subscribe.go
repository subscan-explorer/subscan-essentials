package service

import (
	"context"
	"time"

	"log"

	"github.com/gorilla/websocket"
	"github.com/itering/substrate-api-rpc/rpc"
	ws "github.com/itering/substrate-api-rpc/websocket"
)

const (
	runtimeVersion = iota + 1
	finalizeHeader
)

func (s *Service) Subscribe(ctx context.Context, conn ws.WsConn) {
	var err error

	defer conn.Close()

	subscribeSrv := s.initSubscribeService()
	onceFinHead.Do(func() {
		go subscribeSrv.subscribeFetchBlock(ctx)
	})
	go func() {
		for {
			if !conn.IsConnected() {
				continue
			}
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("read: %s", err)
				continue
			}
			_ = subscribeSrv.parser(message)
		}
	}()

	if err = conn.WriteMessage(websocket.TextMessage, rpc.ChainGetRuntimeVersion(runtimeVersion)); err != nil {
		log.Printf("write: %s", err)
	}
	if err = conn.WriteMessage(websocket.TextMessage, rpc.ChainSubscribeFinalizedHeads(finalizeHeader)); err != nil {
		log.Printf("write: %s", err)
	}

	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Printf("write close: %s", err)
				return
			}
			conn.Close()
			return
		}
	}
}
