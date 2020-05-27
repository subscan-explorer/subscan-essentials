package daemons

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-kratos/kratos/pkg/log"
	"github.com/gorilla/websocket"
	"github.com/itering/subscan/internal/substrate/rpc"
	"github.com/itering/subscan/pkg/recws"
	"github.com/itering/subscan/util"
)

var (
	lockId         sync.Mutex
	substrateRpcId int
	SubscribeConn  *recws.RecConn
	subscribeOnce  sync.Once
)

func Subscribe() {
	var err error

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	SubscribeConn = &recws.RecConn{KeepAliveTimeout: 10 * time.Second}
	SubscribeConn.Dial(util.WSEndPoint, nil)
	defer SubscribeConn.Close()

	subscribeSrv := srv.InitSubscribeService()
	done := make(chan struct{})
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
			subscribeSrv.ParserSubscribe(message)
		}
	}()
	subscribeOnce.Do(func() {
		if err = SubscribeConn.WriteMessage(websocket.TextMessage, rpc.ChainGetRuntimeVersion(1)); err != nil {
			log.Info("write: %s", err)
		}
		if err = SubscribeConn.WriteMessage(websocket.TextMessage, rpc.ChainSubscribeNewHead(101)); err != nil {
			log.Info("write: %s", err)
		}
		if err = SubscribeConn.WriteMessage(websocket.TextMessage, rpc.ChainSubscribeFinalizedHeads(102)); err != nil {
			log.Info("write: %s", err)
		}
	})
	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()
	substrateRpcId = 1
	for {
		select {
		case <-done:
			return
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
func checkHealth() {
	lockId.Lock()
	defer lockId.Unlock()
	substrateRpcId++
	if err := SubscribeConn.WriteMessage(websocket.TextMessage, rpc.SystemHealth(substrateRpcId)); err != nil {
		log.Info("SystemHealth get error: %v", err)
	}
	if substrateRpcId >= 100 {
		substrateRpcId = 2
	}

}
