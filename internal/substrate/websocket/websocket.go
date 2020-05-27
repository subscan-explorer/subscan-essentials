package websocket

import (
	"fmt"
	"github.com/itering/subscan/internal/pkg/recws"
	"github.com/itering/subscan/internal/util"
	"time"
)

func Init() (*PoolConn, error) {
	var err error
	if wsPool == nil {
		factory := func() (*recws.RecConn, error) {
			SubscribeConn := &recws.RecConn{KeepAliveTimeout: 10 * time.Second}
			SubscribeConn.Dial(util.WSEndPoint, nil)
			return SubscribeConn, err
		}
		if wsPool, err = NewChannelPool(1, 10, factory); err != nil {
			fmt.Println("NewChannelPool", err)
		}
	}
	if err != nil {
		return nil, err
	}
	conn, err := wsPool.Get()
	return conn, err
}

func CloseWsConnection() {
	if wsPool != nil {
		wsPool.Close()
	}
}
