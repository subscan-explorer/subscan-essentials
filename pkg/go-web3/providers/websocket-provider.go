/********************************************************************************
   This file is part of go-web3.
   go-web3 is free software: you can redistribute it and/or modify
   it under the terms of the GNU Lesser General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.
   go-web3 is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Lesser General Public License for more details.
   You should have received a copy of the GNU Lesser General Public License
   along with go-web3.  If not, see <http://www.gnu.org/licenses/>.
*********************************************************************************/

/**
 * @file websocket-provider.go
 * @authors:
 *   Reginaldo Costa <regcostajr@gmail.com>
 * @date 2017
 */

package providers

import (
	"math/rand"

	"subscan/pkg/go-web3/constants"

	"golang.org/x/net/websocket"
	"subscan/pkg/go-web3/providers/util"
)

type WebSocketProvider struct {
	address string
	ws      *websocket.Conn
}

func NewWebSocketProvider(address string) *WebSocketProvider {
	provider := new(WebSocketProvider)
	provider.address = address
	return provider
}

func (provider WebSocketProvider) SendRequest(v interface{}, method string, params interface{}) error {

	bodyString := util.JSONRPCObject{Version: "2.0", Method: method, Params: params, ID: rand.Intn(100)}

	if provider.ws == nil {
		ws, err := websocket.Dial(provider.address, "", provider.address)
		if err != nil {
			return err
		}
		provider.ws = ws
	}

	message := []byte(bodyString.AsJsonString())
	_, err := provider.ws.Write(message)
	if err != nil {
		return err
	}

	return websocket.JSON.Receive(provider.ws, v)

}

func (provider WebSocketProvider) Close() error {
	if provider.ws != nil {
		return provider.ws.Close()
	}

	return customerror.WEBSOCKETNOTDENIFIED

}
