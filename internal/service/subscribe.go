package service

import (
	"github.com/go-kratos/kratos/pkg/log"
	"github.com/gorilla/websocket"
	"github.com/itering/substrate-api-rpc/rpc"
	ws "github.com/itering/substrate-api-rpc/websocket"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	runtimeVersion = iota + 1
	newHeader
	finalizeHeader
)

func (s *Service) Subscribe(conn ws.WsConn, interrupt chan os.Signal) {
	var err error

	signal.Notify(interrupt, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

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
				log.Error("read: %s", err)
				continue
			}
			_ = subscribeSrv.parser(message)
		}
	}()

	if err = conn.WriteMessage(websocket.TextMessage, rpc.ChainGetRuntimeVersion(runtimeVersion)); err != nil {
		log.Info("write: %s", err)
	}
	if err = conn.WriteMessage(websocket.TextMessage, rpc.ChainSubscribeNewHead(newHeader)); err != nil {
		log.Info("write: %s", err)
	}
	if err = conn.WriteMessage(websocket.TextMessage, rpc.ChainSubscribeFinalizedHeads(finalizeHeader)); err != nil {
		log.Info("write: %s", err)
	}

	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			if err := conn.WriteMessage(websocket.TextMessage, rpc.SystemHealth(rand.Intn(100)+finalizeHeader)); err != nil {
				log.Info("SystemHealth get error: %v", err)
			}
		case <-interrupt:
			close(done)
			log.Info("interrupt")
			err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Error("write close: %s", err)
				return
			}
			return
		}
	}

}
