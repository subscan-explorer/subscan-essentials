package websocket

import (
	"context"
	"encoding/json"
	"github.com/bilibili/kratos/pkg/log"
	"github.com/garyburd/redigo/redis"
	"subscan-end/internal/service"
	"subscan-end/utiles"
	"time"
)

type Router struct {
	BroadcastConn chan []byte
	srv           *service.Service
}

type SystemMessage struct {
	Topic   string      `json:"topic"`
	Content interface{} `json:"content"`
}

func NewMessageRouter(bc chan []byte, srv *service.Service) *Router {
	r := &Router{BroadcastConn: bc, srv: srv}
	ctx := context.TODO()
	go r.msgSub(ctx)
	return r
}

func (r *Router) msgSub(ctx context.Context) {
	conn := utiles.SubPool.Get()
	psc := redis.PubSubConn{Conn: conn}
	_ = psc.Subscribe(utiles.SubScanChannel)
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			r.dispatch(v.Data)
		case redis.Subscription:
			log.Info("subscribe success")
		case error:
			_ = conn.Close()
			psc = redis.PubSubConn{Conn: utiles.SubPool.Get()}
			_ = psc.Subscribe(utiles.SubScanChannel)
		}
	}
}

func (r *Router) dispatch(data []byte) {
	var publishData SystemMessage
	err := json.Unmarshal(data, &publishData)
	if err != nil {
		log.Error("dispatch Unmarshal error %s", string(data))
	} else {
		newPush := map[string]interface{}{
			"topic":   publishData.Topic,
			"time":    time.Now().Unix(),
			"content": publishData.Content,
		}
		pushBytes, _ := json.Marshal(newPush)
		r.BroadcastConn <- pushBytes
	}
}
