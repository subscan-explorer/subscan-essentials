package service

import (
	"github.com/itering/subscan/lib/substrate/rpc"
	"github.com/itering/subscan/lib/substrate/storageKey"
	"github.com/itering/subscan/pkg/recws"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kratos/kratos/pkg/log"
	"github.com/gorilla/websocket"
	"github.com/itering/subscan/util"
)

var (
	subscribeConn *recws.RecConn
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

	subscribeConn = &recws.RecConn{KeepAliveTimeout: 10 * time.Second}
	subscribeConn.Dial(util.WSEndPoint, nil)

	for {
		if subscribeConn.IsConnected() {
			break
		}
		time.Sleep(subscribeConn.RecIntvlMin)
	}

	defer subscribeConn.Close()

	done := make(chan struct{})

	subscribeSrv := s.InitSubscribeService(done)
	go func() {
		for {
			if !subscribeConn.IsConnected() {
				continue
			}
			_, message, err := subscribeConn.ReadMessage()
			if err != nil {
				log.Error("read: %s", err)
				continue
			}
			log.Info("recv: %s", message)
			subscribeSrv.Parser(message)
		}
	}()

	if err = subscribeConn.WriteMessage(websocket.TextMessage, rpc.ChainGetRuntimeVersion(runtimeVersion)); err != nil {
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
					if err = subscribeConn.WriteMessage(websocket.TextMessage, rpc.ChainSubscribeNewHead(newHeader)); err != nil {
						log.Info("write: %s", err)
					}
				case finalizeHeader:
					if err = subscribeConn.WriteMessage(websocket.TextMessage, rpc.ChainSubscribeFinalizedHeads(finalizeHeader)); err != nil {
						log.Info("write: %s", err)
					}

				case stateChange:
					if err = subscribeConn.WriteMessage(websocket.TextMessage, rpc.StateSubscribeStorage(stateChange, subscribeStorageList)); err != nil {
						log.Info("write: %s", err)
					}
				}
			}
		}
	}

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			checkHealth()
		case <-interrupt:
			close(done)
			log.Info("interrupt")
			err = subscribeConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Error("write close: %s", err)
				return
			}

			return
		}
	}

}
