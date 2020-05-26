package websocket

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/itering/subscan/pkg/recws"
	"sync"
)

var (
	ErrClosed = errors.New("pool is closed")
)

type Pool interface {
	Get() (*PoolConn, error)
	Close()
	Len() int
}

type channelPool struct {
	mu      sync.RWMutex
	conns   chan *recws.RecConn
	factory Factory
}

type Factory func() (*recws.RecConn, error)

func NewChannelPool(initialCap, maxCap int, factory Factory) (Pool, error) {
	if initialCap < 0 || maxCap <= 0 || initialCap > maxCap {
		return nil, errors.New("invalid capacity settings")
	}
	c := &channelPool{
		conns:   make(chan *recws.RecConn, maxCap),
		factory: factory,
	}
	for i := 0; i < initialCap; i++ {
		conn, err := factory()
		if err != nil {
			c.Close()
			return nil, fmt.Errorf("factory is not able to fill the pool: %s", err)
		}
		c.conns <- conn
	}
	return c, nil
}

func (c *channelPool) getConnsAndFactory() (chan *recws.RecConn, Factory) {
	c.mu.RLock()
	conns := c.conns
	factory := c.factory
	c.mu.RUnlock()
	return conns, factory
}

func (c *channelPool) Get() (*PoolConn, error) {
	conns, factory := c.getConnsAndFactory()
	if conns == nil {
		return nil, ErrClosed
	}
	var err error
	select {
	case conn := <-conns:
		if conn == nil || conn.IsConnected() == false {
			conn, err = factory()
			if err != nil {
				return nil, err
			}
		}
		return c.wrapConn(conn), nil
	default:
		conn, err := factory()
		if err != nil {
			return nil, err
		}

		return c.wrapConn(conn), nil
	}
}

func (c *channelPool) put(conn *recws.RecConn) error {
	if conn == nil {
		return errors.New("connection is nil. rejecting")
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.conns == nil {
		conn.Close()
		return nil
	}

	select {
	case c.conns <- conn:
		return nil
	default:
		conn.Close()
		return nil
	}
}

func (c *channelPool) Close() {
	c.mu.Lock()
	conns := c.conns
	c.conns = nil
	c.factory = nil
	c.mu.Unlock()

	if conns == nil {
		return
	}

	close(conns)
	for conn := range conns {
		conn.Close()
	}
}

func (c *channelPool) Len() int {
	conns, _ := c.getConnsAndFactory()
	return len(conns)
}

type PoolConn struct {
	Conn     *recws.RecConn
	mu       sync.RWMutex
	c        *channelPool
	unusable bool
}

// Close() puts the given connects back to the pool instead of closing it.
func (p *PoolConn) Close() error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.unusable {
		if p.Conn != nil {
			p.Conn.Close()
			return nil
		}
		return nil
	}
	return p.c.put(p.Conn)
}

// MarkUnusable() marks the connection not usable any more, to let the pool close it instead of returning it to pool.
func (p *PoolConn) MarkUnusable() {
	p.mu.Lock()
	p.unusable = true
	p.mu.Unlock()
}

// newConn wraps a standard net.Conn to a poolConn net.Conn.
func (c *channelPool) wrapConn(conn *recws.RecConn) *PoolConn {
	p := &PoolConn{c: c}
	p.Conn = conn
	return p
}

var wsPool Pool

func SendWsRequest(c *recws.RecConn, v interface{}, action []byte) (err error) {
	var p *PoolConn
	if c == nil {
		if p, err = Init(); err != nil {
			return
		}
		defer p.Close()
		c = p.Conn
	}
	if err = c.WriteMessage(websocket.TextMessage, action); err != nil {
		if p != nil {
			p.MarkUnusable()
		}
		return fmt.Errorf("websocket send error: %v", err)
	}
	if err = c.ReadJSON(v); err != nil {
		if p != nil {
			p.MarkUnusable()
		}
		return
	}
	return nil
}
