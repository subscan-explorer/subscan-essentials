package substrate

import (
	"fmt"
	"github.com/gorilla/websocket"
	"subscan-end/utiles"
)

var WsConnection *websocket.Conn

func SendWsRequest(v interface{}, action []byte) error {
	c, err := initWebsocket()
	if err != nil {
		return nil
	}
	if err = c.WriteMessage(websocket.TextMessage, action); err != nil {
		return fmt.Errorf("websocket send error: %v", err)
	}
	return c.ReadJSON(v)
}

func initWebsocket() (*websocket.Conn, error) {
	var err error
	if WsConnection == nil {
		WsConnection, _, err = websocket.DefaultDialer.Dial(utiles.ProviderEndPoint, nil)
		if err != nil {
			return nil, err
		}
	}
	return WsConnection, err
}

func CloseWsConnection() {
	if WsConnection != nil {
		_ = WsConnection.Close()
	}
}
