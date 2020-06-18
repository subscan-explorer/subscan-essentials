package service

import (
	"github.com/itering/subscan/internal/pkg/recws"
	"github.com/itering/subscan/internal/substrate/rpc"
	"github.com/itering/subscan/internal/substrate/storageKey"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kratos/kratos/pkg/log"
	"github.com/gorilla/websocket"
	"github.com/itering/subscan/internal/util"
)

var (
	SubscribeConn *recws.RecConn
)

const (
	subscribeTimeoutInterval = 30

	runtimeVersion = iota + 1
	newHeader
	finalizeHeader
	stateChange
)

func (s *Service) Subscribe() {
	var err error

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	SubscribeConn = &recws.RecConn{KeepAliveTimeout: 10 * time.Second}
	SubscribeConn.Dial(util.WSEndPoint, nil)

	for {
		if SubscribeConn.IsConnected() {
			break
		}
		time.Sleep(SubscribeConn.RecIntvlMin)
	}

	defer SubscribeConn.Close()

	done := make(chan struct{})

	subscribeSrv := s.InitSubscribeService(done)
	go func() {
		defer close(done)
		for {
			if !SubscribeConn.IsConnected() {
				continue
			}
			_, message, err := SubscribeConn.ReadMessage()
			if err != nil {
				log.Error("read: %s", err)
				return
			}
			log.Info("recv: %s", message)
			subscribeSrv.Parser(message)
		}
	}()

	if err = SubscribeConn.WriteMessage(websocket.TextMessage, rpc.ChainGetRuntimeVersion(runtimeVersion)); err != nil {
		log.Info("write: %s", err)
	}

	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()

	subscribeStorageList := storageKey.SubscribeStorage()
	checkHealth := func() {
		for _, subscript := range subscriptionIds {
			if time.Now().Unix()-subscript.Latest > subscribeTimeoutInterval {
				switch subscript.Topic {

				case newHeader:
					if err = SubscribeConn.WriteMessage(websocket.TextMessage, rpc.ChainSubscribeNewHead(newHeader)); err != nil {
						log.Info("write: %s", err)
					}
				case finalizeHeader:
					if err = SubscribeConn.WriteMessage(websocket.TextMessage, rpc.ChainSubscribeFinalizedHeads(finalizeHeader)); err != nil {
						log.Info("write: %s", err)
					}

				case stateChange:
					if err = SubscribeConn.WriteMessage(websocket.TextMessage, rpc.StateSubscribeStorage(stateChange, subscribeStorageList)); err != nil {
						log.Info("write: %s", err)
					}
				}
			}
		}
	}

	for {
		select {
		// case <-done:
		// 	return
		case <-ticker.C:
			checkHealth()
		case <-interrupt:
			log.Info("interrupt")
			err = SubscribeConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Error("write close: %s", err)
				return
			}
			return
		}
	}

}
