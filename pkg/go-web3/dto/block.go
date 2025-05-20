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
 * @file block.go
 * @authors:
 *   Jérôme Laurens <jeromelaurens@gmail.com>
 * @date 2017
 */

package dto

type Block struct {
	Difficulty      string             `json:"difficulty"`
	ExtraData       string             `json:"extraData"`
	GasLimit        string             `json:"gasLimit"`
	GasUsed         string             `json:"gasUsed"`
	Author          string             `json:"author,omitempty"`
	Hash            string             `json:"hash"`
	LogsBloom       string             `json:"logs_bloom"`
	Miner           string             `json:"miner"`
	MixHash         string             `json:"mixHash"`
	Nonce           string             `json:"nonce"`
	Number          string             `json:"number"`
	ParentHash      string             `json:"parentHash"`
	ReceiptsRoot    string             `json:"receiptsRoot"`
	Sha3Uncles      string             `json:"sha3Uncles"`
	StateRoot       string             `json:"stateRoot"`
	Size            string             `json:"size"`
	Timestamp       string             `json:"timestamp"`
	SealFields      []string           `json:"sealFields"`
	Uncles          []string           `json:"uncles"`
	TotalDifficulty string             `json:"totalDifficulty"`
	Transactions    []BlockTransaction `json:"transactions"`

	TransactionsRoot string `json:"transactionsRoot"`
	BaseFeePerGas    string `json:"baseFeePerGas"`
}

type BlockTransaction struct {
	Hash                 string `json:"hash"`
	Nonce                string `json:"nonce"`
	BlockHash            string `json:"blockHash"`
	BlockNumber          string `json:"blockNumber"`
	TransactionIndex     string `json:"transactionIndex"`
	From                 string `json:"from"`
	To                   string `json:"to"`
	Input                string `json:"input"`
	Value                string `json:"value"`
	GasPrice             string `json:"gasPrice,omitempty"`
	Gas                  string `json:"gas,omitempty"`
	R                    string `json:"r"`
	S                    string `json:"s"`
	V                    string `json:"v"`
	Creates              string `json:"creates"`
	MaxPriorityFeePerGas string `json:"maxPriorityFeePerGas"`
	MaxFeePerGas         string `json:"maxFeePerGas"`
	Type                 string `json:"type,omitempty"`
}
