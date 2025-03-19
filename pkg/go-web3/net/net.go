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
 * @file net.go
 * @authors:
 *   Reginaldo Costa <regcostajr@gmail.com>
 * @date 2017
 */

package net

import (
	"context"
	"math/big"
	"subscan/pkg/go-web3/dto"
	"subscan/pkg/go-web3/providers"
)

// Net - The Net Module
type Net struct {
	provider providers.ProviderInterface
}

// NewNet - Net Module constructor to set the default provider
func NewNet(provider providers.ProviderInterface) *Net {
	net := new(Net)
	net.provider = provider
	return net
}

// IsListening - Returns true if client is actively listening for network connections.
// Reference: https://github.com/ethereum/wiki/wiki/JSON-RPC#net_listening
// Parameters:
//   - none
//
// Returns:
//   - Boolean - true when listening, otherwise false.
func (net *Net) IsListening(ctx context.Context) (bool, error) {

	pointer := &dto.RequestResult{}
	err := net.provider.SendRequest(ctx, pointer, "net_listening", nil)

	if err != nil {
		return false, err
	}

	return pointer.ToBoolean()

}

// GetPeerCount - Returns number of peers currently connected to the client.
// Reference: https://github.com/ethereum/wiki/wiki/JSON-RPC#net_peercount
// Parameters:
//   - none
//
// Returns:
//   - QUANTITY - integer of the number of connected peers.
func (net *Net) GetPeerCount(ctx context.Context) (*big.Int, error) {

	pointer := &dto.RequestResult{}

	err := net.provider.SendRequest(ctx, pointer, "net_peerCount", nil)

	if err != nil {
		return nil, err
	}

	return pointer.ToBigInt()

}

// GetVersion - Returns the current network id.
// Reference: https://github.com/ethereum/wiki/wiki/JSON-RPC#net_version
// Parameters:
//   - none
//
// Returns:
//   - String - The current network id.
//     "1": Ethereum Mainnet
//     "2": Morden Testnet (deprecated)
//     "3": Ropsten Testnet
//     "4": Rinkeby Testnet
//     "42": Kovan Testnet
func (net *Net) GetVersion(ctx context.Context) (string, error) {

	pointer := &dto.RequestResult{}

	err := net.provider.SendRequest(ctx, pointer, "net_version", nil)

	if err != nil {
		return "", err
	}

	return pointer.ToString()

}
