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
 * @file web3.go
 * @authors:
 *   Reginaldo Costa <regcostajr@gmail.com>
 * @date 2017
 */

package web3

import (
	"context"
	"github.com/itering/subscan/pkg/go-web3/dto"
	"github.com/itering/subscan/pkg/go-web3/eth"
	"github.com/itering/subscan/pkg/go-web3/net"
	"github.com/itering/subscan/pkg/go-web3/personal"
	"github.com/itering/subscan/pkg/go-web3/providers"
	"github.com/itering/subscan/pkg/go-web3/utils"
)

// Coin - Ethereum value unity value
const (
	Coin float64 = 1000000000000000000
)

// Web3 - The Web3 Module
type Web3 struct {
	Provider providers.ProviderInterface
	Eth      *eth.Eth
	Net      *net.Net
	Personal *personal.Personal
	Utils    *utils.Utils
}

// NewWeb3 - Web3 Module constructor to set the default provider, Eth, Net and Personal
func NewWeb3(provider providers.ProviderInterface) *Web3 {
	web3 := new(Web3)
	web3.Provider = provider
	web3.Eth = eth.NewEth(provider)
	web3.Net = net.NewNet(provider)
	web3.Personal = personal.NewPersonal(provider)
	web3.Utils = utils.NewUtils(provider)
	return web3
}

// ClientVersion - Returns the current client version.
// Reference: https://github.com/ethereum/wiki/wiki/JSON-RPC#web3_clientversion
// Parameters:
//   - none
//
// Returns:
//   - String - The current client version
func (web Web3) ClientVersion(ctx context.Context) (string, error) {

	pointer := &dto.RequestResult{}

	err := web.Provider.SendRequest(ctx, pointer, "web3_clientVersion", nil)

	if err != nil {
		return "", err
	}

	return pointer.ToString()

}
