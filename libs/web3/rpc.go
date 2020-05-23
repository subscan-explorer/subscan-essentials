package web3

import (
	"encoding/json"
	"fmt"

	"github.com/itering/subscan/util"
)

// curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getTransactionByBlockHashAndIndex","params":["0xccb81797924fc669d1de5e7bbdd38f89d87157c799221051f264a026fb924a3a", "0x38"],"id":1}'

const endpoint = "https://ropsten.infura.io/v3/1bb85682d6494e219803bab49a4813dc"

type ReqBody struct {
	JSONRPC string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	ID      int      `json:"id"`
	Params  []string `json:"params"`
}

type Transaction struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  struct {
		BlockHash   string `json:"blockHash"`
		BlockNumber string `json:"blockNumber"`
		From        string `json:"from"`
		Gas         string `json:"gas"`
		GasPrice    string `json:"gasPrice"`
		Hash        string `json:"hash"`
	} `json:"result"`
}

func EthGetTransactionByBlockHashAndIndex(blockHash string, index int) string {
	r := ReqBody{
		JSONRPC: "2.0",
		Method:  "eth_getTransactionByBlockHashAndIndex",
		ID:      1,
		Params:  []string{blockHash, fmt.Sprintf("0x%x", index)},
	}
	j, _ := json.Marshal(r)
	b, err := util.PostWithJson(j, endpoint)
	if b == nil || err != nil {
		return ""
	}
	var transaction Transaction
	err = json.Unmarshal(b, &transaction)
	if err != nil {
		return ""
	}
	return transaction.Result.Hash

}
