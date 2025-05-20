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
 * @file complex-int.go
 * @authors:
 *   Reginaldo Costa <regcostajr@gmail.com>
 * @date 2017
 */

package types

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

type ComplexIntParameter int64

func (s ComplexIntParameter) ToHex() string {

	return fmt.Sprintf("0x%x", s)

}

type ComplexIntResponse string

func (s ComplexIntResponse) ToUInt64() uint64 {

	sResult, _ := strconv.ParseUint(string(s), 16, 64)
	return sResult

}

func (s ComplexIntResponse) ToInt64() int64 {

	big, _ := new(big.Int).SetString(strings.TrimPrefix(string(s), "0x"), 16)
	return big.Int64()

}

func (s ComplexIntResponse) ToBigInt() *big.Int {
	big, _ := new(big.Int).SetString(string(s), 16)
	return big
}
