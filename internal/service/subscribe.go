package service

import (
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/itering/substrate-api-rpc/storageKey"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kratos/kratos/pkg/log"
	"github.com/gorilla/websocket"
	"github.com/itering/subscan/util"
)

var (
	TotalIssuance storageKey.StorageKey
)

const (
	subscribeTimeoutInterval = 30

	runtimeVersion = iota
	newHeader
	finalizeHeader
	stateChange
)

func subscribeStorage() []string {
	TotalIssuance = storageKey.EncodeStorageKey("Balances", "TotalIssuance")
	return []string{util.AddHex(TotalIssuance.EncodeKey)}
}

type WsConn interface {
	Dial(urlStr string, reqHeader http.Header)
	IsConnected() bool
	Close()
	WriteMessage(messageType int, data []byte) error
	ReadMessage() (messageType int, message []byte, err error)
}

func (s *Service) Subscribe(conn WsConn, interrupt chan os.Signal) {
	var err error

	signal.Notify(interrupt, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	conn.Dial(util.WSEndPoint, nil)

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

	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()

	subscribeStorageList := subscribeStorage()
	checkHealth := func() {
		for _, subscript := range subscriptionIds {
			if time.Now().Unix()-subscript.Latest > subscribeTimeoutInterval {
				switch subscript.Topic {

				case ChainNewHead:
					if err = conn.WriteMessage(websocket.TextMessage, rpc.ChainSubscribeNewHead(newHeader)); err != nil {
						log.Info("write: %s", err)
					}
				case ChainFinalizedHead:
					if err = conn.WriteMessage(websocket.TextMessage, rpc.ChainSubscribeFinalizedHeads(finalizeHeader)); err != nil {
						log.Info("write: %s", err)
					}

				case StateStorage:
					if err = conn.WriteMessage(websocket.TextMessage, rpc.StateSubscribeStorage(stateChange, subscribeStorageList)); err != nil {
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
			err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Error("write close: %s", err)
				return
			}

			return
		}
	}

}
