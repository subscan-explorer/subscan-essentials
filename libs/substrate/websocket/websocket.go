package websocket

import (
	"fmt"
	"time"

	"github.com/itering/subscan/pkg/recws"
	"github.com/itering/subscan/util"
)

func Init() (*PoolConn, error) {
	var err error
	if wsPool == nil {
		factory := func() (*recws.RecConn, error) {
			SubscribeConn := &recws.RecConn{KeepAliveTimeout: 10 * time.Second}
			SubscribeConn.Dial(util.WsEndpointCache(), nil)
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
