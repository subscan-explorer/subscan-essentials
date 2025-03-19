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
 * @file transaction.go
 * @authors:
 *   Reginaldo Costa <regcostajr@gmail.com>
 * @date 2017
 */

package dto

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"subscan/pkg/go-web3/complex/types"
)

// TransactionParameters GO transaction to make more easy controll the parameters
type TransactionParameters struct {
	From     string
	To       string
	Nonce    *big.Int
	Gas      *big.Int
	GasPrice *big.Int
	Value    *big.Int
	Data     types.ComplexString
}

// RequestTransactionParameters JSON
type RequestTransactionParameters struct {
	From     string `json:"from,omitempty"`
	To       string `json:"to,omitempty"`
	Nonce    string `json:"nonce,omitempty"`
	Gas      string `json:"gas,omitempty"`
	GasPrice string `json:"gasPrice,omitempty"`
	Value    string `json:"value,omitempty"`
	Data     string `json:"data,omitempty"`
}

// Transform the GO transactions parameters to json style
func (params *TransactionParameters) Transform() *RequestTransactionParameters {
	request := new(RequestTransactionParameters)
	request.From = params.From
	if params.To != "" {
		request.To = params.To
	}
	if params.Nonce != nil {
		request.Nonce = "0x" + params.Nonce.Text(16)
	}
	if params.Gas != nil {
		request.Gas = "0x" + params.Gas.Text(16)
	}
	if params.GasPrice != nil {
		request.GasPrice = "0x" + params.GasPrice.Text(16)
	}
	if params.Value != nil {
		request.Value = "0x" + params.Value.Text(16)
	}
	if params.Data != "" {
		request.Data = params.Data.ToHex()
	}
	return request
}

type SignTransactionResponse struct {
	Raw         types.ComplexString     `json:"raw"`
	Transaction SignedTransactionParams `json:"tx"`
}

type TracerTransactionResponse struct {
	From    string                      `json:"from"`
	Gas     string                      `json:"gas"`
	GasUsed string                      `json:"gasUsed"`
	Type    string                      `json:"type"`
	To      string                      `json:"to"`
	Output  *string                     `json:"output,omitempty"`
	Input   string                      `json:"input"`
	Value   string                      `json:"value,omitempty"`
	Error   *string                     `json:"error,omitempty"`
	Calls   []TracerTransactionResponse `json:"calls,omitempty"`
}

type SignedTransactionParams struct {
	Gas      *big.Int `json:"gas"`
	GasPrice *big.Int `json:"gasPrice"`
	Hash     string   `json:"hash"`
	Input    string   `json:"input"`
	Nonce    *big.Int `json:"nonce"`
	S        string   `json:"s"`
	R        string   `json:"r"`
	V        *big.Int `json:"v"`
	To       string   `json:"to"`
	Value    *big.Int `json:"value"`
}

type TransactionResponse struct {
	Hash                 string              `json:"hash"`
	Nonce                *big.Int            `json:"nonce"`
	BlockHash            string              `json:"blockHash"`
	BlockNumber          *big.Int            `json:"blockNumber"`
	TransactionIndex     *big.Int            `json:"transactionIndex"`
	From                 string              `json:"from"`
	To                   string              `json:"to"`
	Input                string              `json:"input"`
	Value                *big.Int            `json:"value"`
	GasPrice             *big.Int            `json:"gasPrice,omitempty"`
	Gas                  *big.Int            `json:"gas,omitempty"`
	Data                 types.ComplexString `json:"data,omitempty"`
	R                    string              `json:"r"`
	S                    string              `json:"s"`
	V                    string              `json:"v"`
	Creates              string              `json:"creates"`
	MaxPriorityFeePerGas string              `json:"maxPriorityFeePerGas"`
	MaxFeePerGas         string              `json:"maxFeePerGas"`
}

type TransactionReceipt struct {
	TransactionHash   string            `json:"transactionHash"`
	TransactionIndex  *big.Int          `json:"transactionIndex"`
	BlockHash         string            `json:"blockHash"`
	BlockNumber       *big.Int          `json:"blockNumber"`
	From              string            `json:"from"`
	To                string            `json:"to"`
	CumulativeGasUsed *big.Int          `json:"cumulativeGasUsed"`
	EffectiveGasPrice *big.Int          `json:"effectiveGasPrice,omitempty"`
	GasUsed           *big.Int          `json:"gasUsed"`
	ContractAddress   string            `json:"contractAddress"`
	Logs              []TransactionLogs `json:"logs"`
	LogsBloom         string            `json:"logsBloom"`
	Root              string            `json:"string"`
	Status            bool              `json:"status"`
}

type TransactionLogs struct {
	Address          string   `json:"address"`
	Topics           []string `json:"topics"`
	Data             string   `json:"data"`
	BlockNumber      *big.Int `json:"blockNumber"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex *big.Int `json:"transactionIndex"`
	BlockHash        string   `json:"blockHash"`
	LogIndex         *big.Int `json:"logIndex"`
	Removed          bool     `json:"removed"`
}

func (t *TransactionResponse) UnmarshalJSON(data []byte) error {
	type Alias TransactionResponse
	temp := &struct {
		Nonce            string `json:"nonce"`
		BlockNumber      string `json:"blockNumber,omitempty"`
		TransactionIndex string `json:"transactionIndex"`
		Value            string `json:"value"`
		GasPrice         string `json:"gasPrice,omitempty"`
		Gas              string `json:"gas,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	nonce, success := big.NewInt(0).SetString(temp.Nonce[2:], 16)

	if !success {
		return errors.New(fmt.Sprintf("Error converting nonce %s to BigInt", temp.Nonce))
	}

	if len(temp.BlockNumber) > 0 {
		blockNum, _ := big.NewInt(0).SetString(temp.BlockNumber[2:], 16)
		t.BlockNumber = blockNum
	}

	if len(temp.TransactionIndex) > 0 {
		txIndex, _ := big.NewInt(0).SetString(temp.TransactionIndex[2:], 16)
		t.TransactionIndex = txIndex
	}

	gas, success := big.NewInt(0).SetString(temp.Gas[2:], 16)

	if !success {
		return errors.New(fmt.Sprintf("Error converting Gas  %s to BigInt", temp.Gas))
	}

	gasPrice, success := big.NewInt(0).SetString(temp.GasPrice[2:], 16)

	if !success {
		return errors.New(fmt.Sprintf("Error converting GasPrice %s to BigInt", temp.GasPrice))
	}

	value, success := big.NewInt(0).SetString(temp.Value[2:], 16)

	if !success {
		return errors.New(fmt.Sprintf("Error converting Value %s to BigInt", temp.Value))
	}

	t.Nonce = nonce
	t.Gas = gas
	t.GasPrice = gasPrice
	t.Value = value
	t.R = temp.R
	t.S = temp.S
	t.V = temp.V

	return nil
}

func (r *TransactionLogs) UnmarshalJSON(data []byte) error {
	type Alias TransactionLogs

	log := &struct {
		TransactionIndex string `json:"transactionIndex"`
		BlockNumber      string `json:"blockNumber"`
		LogIndex         string `json:"logIndex"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &log); err != nil {
		return err
	}

	blockNumLog, success := big.NewInt(0).SetString(log.BlockNumber[2:], 16)

	if !success {
		return errors.New(fmt.Sprintf("Error converting %s to BigInt", log.BlockNumber))
	}

	txIndexLogs, success := big.NewInt(0).SetString(log.TransactionIndex[2:], 16)

	if !success {
		return errors.New(fmt.Sprintf("Error converting %s to BigInt", log.TransactionIndex))
	}

	logIndex, success := big.NewInt(0).SetString(log.LogIndex[2:], 16)

	if !success {
		return errors.New(fmt.Sprintf("Error converting %s to BigInt", log.LogIndex))
	}

	r.BlockNumber = blockNumLog
	r.TransactionIndex = txIndexLogs
	r.LogIndex = logIndex
	return nil

}

func (r *TransactionReceipt) UnmarshalJSON(data []byte) error {
	type Alias TransactionReceipt

	temp := &struct {
		TransactionIndex  string `json:"transactionIndex"`
		BlockNumber       string `json:"blockNumber"`
		CumulativeGasUsed string `json:"cumulativeGasUsed"`
		EffectiveGasPrice string `json:"effectiveGasPrice,omitempty"`
		GasUsed           string `json:"gasUsed"`
		Status            string `json:"status"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	blockNum, success := big.NewInt(0).SetString(temp.BlockNumber[2:], 16)

	if !success {
		return errors.New(fmt.Sprintf("Error converting %s to BigInt", temp.BlockNumber))
	}

	txIndex, success := big.NewInt(0).SetString(temp.TransactionIndex[2:], 16)

	if !success {
		return errors.New(fmt.Sprintf("Error converting %s to BigInt", temp.TransactionIndex))
	}

	gasUsed, success := big.NewInt(0).SetString(temp.GasUsed[2:], 16)

	if !success {
		return errors.New(fmt.Sprintf("Error converting %s to BigInt", temp.GasUsed))
	}

	cumulativeGas, success := big.NewInt(0).SetString(temp.CumulativeGasUsed[2:], 16)

	if !success {
		return errors.New(fmt.Sprintf("Error converting %s to BigInt", temp.CumulativeGasUsed))
	}
	var effectiveGasPrice = big.NewInt(0)
	if temp.EffectiveGasPrice != "" {
		effectiveGasPrice, success = big.NewInt(0).SetString(strings.TrimPrefix(temp.EffectiveGasPrice, "0x"), 16)
		if !success {
			return errors.New(fmt.Sprintf("Error converting EffectiveGasPrice %s to BigInt", temp.EffectiveGasPrice))
		}
	}

	status, success := big.NewInt(0).SetString(temp.Status[2:], 16)
	if !success {
		return errors.New(fmt.Sprintf("Error converting Status %s to BigInt", temp.Status))
	}

	r.TransactionIndex = txIndex
	r.BlockNumber = blockNum
	r.CumulativeGasUsed = cumulativeGas
	r.EffectiveGasPrice = effectiveGasPrice
	r.GasUsed = gasUsed
	r.Status = false
	if status.Cmp(big.NewInt(1)) == 0 {
		r.Status = true
	}

	return nil
}

func (sp *SignedTransactionParams) UnmarshalJSON(data []byte) error {
	type Alias SignedTransactionParams

	temp := &struct {
		Gas      string `json:"gas"`
		GasPrice string `json:"gasPrice"`
		Nonce    string `json:"nonce"`
		V        string `json:"v"`
		Value    string `json:"value"`
		*Alias
	}{
		Alias: (*Alias)(sp),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	gas, success := big.NewInt(0).SetString(temp.Gas[2:], 16)

	if !success {
		return errors.New(fmt.Sprintf("Error converting %s to BigInt", temp.Gas))
	}

	gasPrice, success := big.NewInt(0).SetString(temp.GasPrice[2:], 16)

	if !success {
		return errors.New(fmt.Sprintf("Error converting %s to BigInt", temp.GasPrice))
	}

	nonce, success := big.NewInt(0).SetString(temp.Nonce[2:], 16)

	if !success {
		return errors.New(fmt.Sprintf("Error converting %s to BigInt", temp.Nonce))
	}

	v, success := big.NewInt(0).SetString(temp.V[2:], 16)

	if !success {
		return errors.New(fmt.Sprintf("Error converting %s to BigInt", temp.V))
	}

	val, success := big.NewInt(0).SetString(temp.Value[2:], 16)
	if !success {
		return errors.New(fmt.Sprintf("Error converting %s to BigInt", temp.Value))
	}

	sp.Gas = gas
	sp.GasPrice = gasPrice
	sp.Nonce = nonce
	sp.V = v
	sp.Value = val

	return nil
}
