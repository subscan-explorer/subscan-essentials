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
 * @file contract.go
 * @authors:
 *   Reginaldo Costa <regcostajr@gmail.com>
 * @date 2018
 */

package eth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/huandu/xstrings"
	"math"
	"strings"
	"subscan/pkg/go-web3/complex/types"
	"subscan/pkg/go-web3/dto"

	"math/big"
	"subscan/pkg/go-web3/utils"
)

// Contract ...
type Contract struct {
	super     *Eth
	abi       string
	functions map[string][]string
}

// NewContract - Contract abstraction
func (eth *Eth) NewContract(abi string) (*Contract, error) {

	contract := new(Contract)
	var mockInterface interface{}

	err := json.Unmarshal([]byte(abi), &mockInterface)

	if err != nil {
		return nil, err
	}

	jsonInterface := mockInterface.([]interface{})
	contract.functions = make(map[string][]string)
	for index := 0; index < len(jsonInterface); index++ {
		function := jsonInterface[index].(map[string]interface{})

		if function["type"] == "constructor" || function["type"] == "fallback" {
			function["name"] = function["type"]
		}

		functionName := function["name"].(string)
		contract.functions[functionName] = make([]string, 0)

		if function["inputs"] == nil {
			continue
		}

		inputs := function["inputs"].([]interface{})
		for paramIndex := 0; paramIndex < len(inputs); paramIndex++ {
			params := inputs[paramIndex].(map[string]interface{})
			contract.functions[functionName] = append(contract.functions[functionName], params["type"].(string))
		}

	}

	contract.abi = abi
	contract.super = eth

	return contract, nil
}

// prepareTransaction ...
func (contract *Contract) prepareTransaction(transaction *dto.TransactionParameters, functionName string, args []interface{}) (*dto.TransactionParameters, error) {

	function, ok := contract.functions[functionName]
	if !ok {
		return nil, errors.New("Function not finded on passed abi")
	}

	fullFunction := functionName + "("

	comma := ""
	for arg := range function {
		fullFunction += comma + function[arg]
		comma = ","
	}

	fullFunction += ")"

	utils := utils.NewUtils(contract.super.provider)
	sha3Function, err := utils.Sha3(types.ComplexString(fullFunction))

	if err != nil {
		return nil, err
	}

	var data string

	for index := 0; index < len(function); index++ {
		currentData, err := contract.getHexValue(function[index], args[index], function...)

		if err != nil {
			return nil, err
		}

		data += currentData
	}

	transaction.Data = types.ComplexString(sha3Function[0:10] + data)

	return transaction, nil

}

func (contract *Contract) Call(ctx context.Context, transaction *dto.TransactionParameters, functionName string, args ...interface{}) (*dto.RequestResult, error) {

	transaction, err := contract.prepareTransaction(transaction, functionName, args)

	if err != nil {
		return nil, err
	}

	return contract.super.Call(ctx, transaction)

}

func (contract *Contract) Send(ctx context.Context, transaction *dto.TransactionParameters, functionName string, args ...interface{}) (string, error) {

	transaction, err := contract.prepareTransaction(transaction, functionName, args)

	if err != nil {
		return "", err
	}

	return contract.super.SendTransaction(ctx, transaction)

}

func (contract *Contract) Deploy(ctx context.Context, transaction *dto.TransactionParameters, bytecode string, args ...interface{}) (string, error) {

	constructor := contract.functions["constructor"]

	for index := 0; index < len(constructor); index++ {
		tmpBytes, err := contract.getHexValue(constructor[index], args[index])

		if err != nil {
			return "", err
		}

		bytecode += tmpBytes
	}

	transaction.Data = types.ComplexString(bytecode)

	return contract.super.SendTransaction(ctx, transaction)

}

func (contract *Contract) getHexValue(inputType string, value interface{}, function ...string) (string, error) {

	var data string

	if inputType == "address[]" || inputType == "uint256[]" {
		return value.(string), nil
	}

	if strings.HasPrefix(inputType, "int") ||
		strings.HasPrefix(inputType, "uint") ||
		strings.HasPrefix(inputType, "fixed") ||
		strings.HasPrefix(inputType, "ufixed") {

		bigVal := value.(*big.Int)

		// Checking that the string actually is the correct inputType
		if strings.Contains(inputType, "128") {
			// 128 bit
			if bigVal.BitLen() > 128 {
				return "", errors.New(fmt.Sprintf("Input type %s not met", inputType))
			}
		} else if strings.Contains(inputType, "256") {
			// 256 bit
			if bigVal.BitLen() > 256 {
				return "", errors.New(fmt.Sprintf("Input type %s not met", inputType))
			}
		}

		data += fmt.Sprintf("%064s", fmt.Sprintf("%x", bigVal))
	}

	if strings.Compare("address", inputType) == 0 {
		data += fmt.Sprintf("%064s", strings.TrimPrefix(value.(string), "0x"))
	}

	if strings.Compare("string", inputType) == 0 {
		data += fmt.Sprintf("%064s", fmt.Sprintf("%x", value.(string)))
	}

	if strings.HasPrefix(inputType, "address[") {
		addressSlice := value.([]string)
		for _, addr := range addressSlice {
			data += fmt.Sprintf("%064s", addr[2:])
		}
	}

	if strings.HasPrefix(inputType, "bytes") {
		if strings.Contains(inputType, "32") {
			data += xstrings.RightJustify(value.(string), 64, "0")
		} else if strings.EqualFold(inputType, "bytes") {
			data += fmt.Sprintf("%064s", fmt.Sprintf("%x", len(function)*32))
			data += fmt.Sprintf("%064s", fmt.Sprintf("%x", int(math.Ceil(float64(len(value.(string)))/2))))
			offset := len(value.(string)) % 64
			valueString := value.(string)

			data += valueString[0:len(valueString)-offset] + xstrings.LeftJustify(valueString[len(valueString)-offset:], 64, "0")
		} else {
			data += xstrings.LeftJustify(strings.TrimPrefix(value.(string), "0x"), 64, "0")
		}

	}

	return data, nil

}
