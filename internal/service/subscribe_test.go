package service

import (
	"github.com/gorilla/websocket"
	"net/http"
	"os"
	"syscall"
	"testing"
)

type TestConn struct {
	*websocket.Conn
	Connected bool
}

func (t *TestConn) Dial(urlStr string, reqHeader http.Header) {
	conn, _, err := websocket.DefaultDialer.Dial(urlStr, nil)
	t.Conn = conn
	if err != nil {
		panic(err)
	}
	t.Connected = true
}

func (t *TestConn) IsConnected() bool {
	return t.Connected
}

func (t *TestConn) Close() {
	t.Connected = false
	if t.Conn != nil {
		t.Conn.Close()
	}
}

func (t *TestConn) ReadMessage() (messageType int, message []byte, err error) {
	_, message, _ = t.Conn.ReadMessage()
	_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	return
}

func TestService_Subscribe(t *testing.T) {
	tc := TestConn{}
	interrupt := make(chan os.Signal, 1)
	testSrv.Subscribe(&tc, interrupt)
}
