package service

import (
	"math/rand"
	"time"

	"log"

	"github.com/gorilla/websocket"
	"github.com/itering/substrate-api-rpc/rpc"
	ws "github.com/itering/substrate-api-rpc/websocket"
)

const (
	runtimeVersion = iota + 1
	newHeader
	finalizeHeader
)

func (s *Service) Subscribe(conn ws.WsConn, stop chan struct{}) {
	var err error

	defer conn.Close()

	done := make(chan struct{})

	subscribeSrv := s.initSubscribeService(done)
	go func() {
		defer close(done)
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
	if err = conn.WriteMessage(websocket.TextMessage, rpc.ChainSubscribeNewHead(newHeader)); err != nil {
		log.Printf("write: %s", err)
	}
	if err = conn.WriteMessage(websocket.TextMessage, rpc.ChainSubscribeFinalizedHeads(finalizeHeader)); err != nil {
		log.Printf("write: %s", err)
	}

	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			if err := conn.WriteMessage(websocket.TextMessage, rpc.SystemHealth(rand.Intn(100)+finalizeHeader)); err != nil {
				log.Printf("SystemHealth get error: %v", err)
			}
		case <-stop:
			close(done)
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
