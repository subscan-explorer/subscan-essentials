package service

import (
	"github.com/gorilla/websocket"
	"net/http"
	"os"
	"strings"
	"syscall"
	"testing"
)

type TestConn struct {
	*websocket.Conn
	Connected bool
}

func (t *TestConn) MarkUnusable() {}

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
	return
}

func (t *TestConn) WriteMessage(messageType int, data []byte) error {
	if strings.EqualFold(string(data), `{"id":4,"method":"state_subscribeStorage","params":[["0x481e203dcea218263e3a96ca9e4b193857c875e4cff74148e4628f264b974c80"]],"jsonrpc":"2.0"}`) {
		_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}
	return nil
}

func TestService_Subscribe(t *testing.T) {
	tc := TestConn{}
	interrupt := make(chan os.Signal, 1)
	testSrv.Subscribe(&tc, interrupt)
}
