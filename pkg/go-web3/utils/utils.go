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
 * @file utils.go
 * @authors:
 *   Reginaldo Costa <regcostajr@gmail.com>
 * @date 2017
 */

package utils

import (
	"encoding/hex"
	"github.com/itering/subscan/pkg/go-web3/complex/types"
	"github.com/itering/subscan/pkg/go-web3/providers"
	"strings"
)

// Utils - The Utils Module
type Utils struct {
	provider providers.ProviderInterface
}

// NewUtils - Utils Module constructor to set the default provider
func NewUtils(provider providers.ProviderInterface) *Utils {
	utils := new(Utils)
	utils.provider = provider
	return utils
}

// Sha3 - Returns Keccak-256 (not the standardized SHA3-256) of the given data.
// Reference: https://github.com/ethereum/wiki/wiki/JSON-RPC#web3_sha3
//   - DATA - the data to convert into a SHA3 hash
//
// Returns:
//   - DATA - The SHA3 result of the given string.
func (utils *Utils) Sha3(data types.ComplexString) (string, error) {
	return utils.BytesToHex(Keccak256(utils.HexToBytes(data.ToHex()))), nil
}

func (utils *Utils) HexToBytes(s string) []byte {
	s = strings.TrimPrefix(s, "0x")
	c := make([]byte, hex.DecodedLen(len(s)))
	_, _ = hex.Decode(c, []byte(s))
	return c
}

func (utils *Utils) BytesToHex(b []byte) string {
	c := make([]byte, hex.EncodedLen(len(b)))
	hex.Encode(c, b)
	return "0x" + string(c)
}
