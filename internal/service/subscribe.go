package service

import (
	"context"
	"fmt"
	"github.com/itering/subscan/util"
	"sync"
	"time"

	"log"

	"github.com/gorilla/websocket"
	"github.com/itering/substrate-api-rpc/rpc"
)

const (
	runtimeVersion = iota + 1
	finalizeHeader
	newHeader
)

func subscribeFromChain() (err error) {
	if conn != nil {
		defer func() {
			if err == nil {
				util.Logger().Info("subscribe from chain success!")
			} else {
				util.Logger().Error(fmt.Errorf("subscribe from chain failed: %v", err))
			}
		}()
		if err = conn.WriteMessage(websocket.TextMessage, rpc.ChainGetRuntimeVersion(runtimeVersion)); err != nil {
			log.Printf("write: %s", err)
		}
		if err = conn.WriteMessage(websocket.TextMessage, rpc.ChainSubscribeFinalizedHeads(finalizeHeader)); err != nil {
			log.Printf("write: %s", err)
		}
		if err = conn.WriteMessage(websocket.TextMessage, rpc.ChainSubscribeNewHead(newHeader)); err != nil {
			log.Printf("write: %s", err)
		}
	}
	return
}

var (
	conn      *websocket.Conn
	connMutex sync.RWMutex
)

func getConn() *websocket.Conn {
	connMutex.RLock()
	defer connMutex.RUnlock()
	return conn
}

func setConn(newConn *websocket.Conn) {
	connMutex.Lock()
	defer connMutex.Unlock()
	if conn != nil {
		safeClose(conn)
	}
	conn = newConn
}

func safeClose(c *websocket.Conn) {
	if c != nil {
		msg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "bye")
		_ = c.WriteControl(websocket.CloseMessage, msg, time.Now().Add(time.Second))
		_ = c.Close()
	}
}

func reSubscribeFromChain() {
	for {
		setConn(nil) // close old connection

		newConn, _, err := websocket.DefaultDialer.Dial(util.WSEndPoint, nil)
		if err != nil {
			util.Logger().Error(fmt.Errorf("dial error: %v", err))
			time.Sleep(time.Second * 2)
			continue
		}

		setConn(newConn)

		if err = subscribeFromChain(); err != nil {
			util.Logger().Error(fmt.Errorf("subscribe error: %v", err))
			continue
		}
		break
	}
}

func (s *Service) Subscribe(ctx context.Context) {
	reSubscribeFromChain()
	defer safeClose(getConn())

	subscribeSrv := s.initSubscribeService()
	onceFinHead.Do(func() {
		go subscribeSrv.subscribeFetchBlock(ctx)
	})

	go func() {
		for {
			c := getConn()
			if c == nil {
				continue
			}

			_, message, err := c.ReadMessage()
			if err != nil {
				time.Sleep(time.Second * 5)
				util.Logger().Error(fmt.Errorf("read error: %v", err))
				reSubscribeFromChain()
				continue
			}
			_ = subscribeSrv.parser(message)
		}
	}()

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c := getConn()
			if c == nil {
				continue
			}
			_ = c.SetWriteDeadline(time.Now().Add(5 * time.Second))
			if err := c.WriteMessage(websocket.PingMessage, nil); err != nil {
				util.Logger().Error(fmt.Errorf("ping error: %v", err))
			}
		}
	}
}
