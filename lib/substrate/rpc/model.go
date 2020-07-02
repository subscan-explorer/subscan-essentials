package rpc

import (
	"encoding/json"
	"errors"
	"github.com/shopspring/decimal"
)

var nilErr = errors.New("nil result")

type JsonRpcResult struct {
	Id      int         `json:"id,omitempty"`
	JsonRpc string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Params  *SubParams  `json:"params,omitempty"`
	Method  string      `json:"method,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}
type JsonRpcParams struct {
	Id      int         `json:"id,omitempty"`
	JsonRpc string      `json:"jsonrpc"`
	Params  interface{} `json:"params"`
	Method  string      `json:"method"`
	Error   *Error      `json:"error,omitempty"`
}

type Properties struct {
	Ss58Format    int    `json:"ss58Format"`
	TokenDecimals int    `json:"tokenDecimals"`
	TokenSymbol   string `json:"tokenSymbol"`
}

type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type SubParams struct {
	Result interface{} `json:"result"`
}

type HealthResult struct {
	IsSyncing       bool `json:"is_syncing"`
	Peers           int  `json:"peers"`
	ShouldHavePeers bool `json:"shouldHavePeers"`
}

type RuntimeVersion struct {
	Apis             [][]interface{} `json:"apis"`
	AuthoringVersion int             `json:"authoringVersion"`
	ImplName         string          `json:"implName"`
	ImplVersion      int             `json:"implVersion"`
	SpecName         string          `json:"specName"`
	SpecVersion      int             `json:"specVersion"`
}

type SystemTokenResult struct {
	TokenDecimals int    `json:"tokenDecimals"`
	TokenSymbol   string `json:"tokenSymbol"`
}

type ChainNewHeadResult struct {
	ExtrinsicsRoot string          `json:"extrinsicsRoot"`
	Number         string          `json:"number"`
	ParentHash     string          `json:"parentHash"`
	StateRoot      string          `json:"stateRoot"`
	Digest         ChainNewHeadLog `json:"digest"`
}

type PaymentQueryInfo struct {
	Class      string          `json:"class"`
	PartialFee decimal.Decimal `json:"partialFee"`
	Weight     int64           `json:"weight"`
}

type ChainNewHeadLog struct {
	Logs []string `json:"logs"`
}

type StateStorageResult struct {
	Block   string     `json:"block"`
	Changes [][]string `json:"changes"`
}

type BlockResult struct {
	Block         Block  `json:"block"`
	Justification string `json:"justification"`
}

type Block struct {
	Extrinsics []string           `json:"extrinsics"`
	Header     ChainNewHeadResult `json:"header"`
}

func (p *JsonRpcResult) ToString() (string, error) {
	if p.checkErr() != nil {
		return "", p.checkErr()
	}
	if p.Result == nil {
		return "", nil
	}
	return p.Result.(string), nil
}

func (p *JsonRpcResult) ToInterfaces() ([]interface{}, error) {
	if p.checkErr() != nil {
		return nil, p.checkErr()
	}
	if p.Result == nil {
		return nil, nil
	}
	return p.Result.([]interface{}), nil
}

func (p *JsonRpcResult) ToInt() uint64 {
	if p.checkErr() != nil {
		return 0
	}
	return p.Result.(uint64)
}

func (p *JsonRpcResult) ToFloat64() float64 {
	if p.checkErr() != nil {
		return 0
	}
	return p.Result.(float64)
}

func (p *JsonRpcResult) ToRuntimeVersion() *RuntimeVersion {
	if p.checkErr() != nil {
		return nil
	}
	result := (p).Result.(map[string]interface{})
	if len(result) == 0 {
		return nil
	}
	v := &RuntimeVersion{}
	marshal, err := json.Marshal(result)
	if err != nil {
		return nil
	}
	_ = json.Unmarshal([]byte(marshal), v)
	return v
}

func (p *JsonRpcResult) ToAnyThing(r interface{}) error {
	if p.checkErr() != nil {
		return p.checkErr()
	}
	result := (p).Result.(map[string]interface{})
	if len(result) == 0 {
		return nilErr
	}
	marshal, err := json.Marshal(result)
	if err != nil {
		return nilErr
	}
	_ = json.Unmarshal([]byte(marshal), r)
	return nil
}

func (p *JsonRpcResult) ToStorage() (*StateStorageResult, int64) {
	if p.Params == nil {
		return nil, 0
	}
	result := (p).Params.Result.(map[string]interface{})
	if len(result) == 0 {
		return nil, 0
	}
	v := &StateStorageResult{}
	marshal, err := json.Marshal(result)
	if err != nil {
		return nil, 0
	}
	_ = json.Unmarshal([]byte(marshal), v)
	return v, 0
}

func (p *JsonRpcResult) ToNewHead() *ChainNewHeadResult {
	if p.checkErr() != nil {
		return nil
	}
	if p.Params == nil {
		return nil
	}
	result := (p).Params.Result.(map[string]interface{})
	if len(result) == 0 {
		return nil
	}
	v := &ChainNewHeadResult{}
	marshal, err := json.Marshal(result)
	if err != nil {
		return nil
	}
	_ = json.Unmarshal([]byte(marshal), v)
	return v
}

func (p *JsonRpcResult) ToSysHealth() *HealthResult {
	if p.checkErr() != nil {
		return nil
	}
	result := (p).Result.(map[string]interface{})
	if len(result) == 0 {
		return nil
	}
	v := &HealthResult{}
	marshal, err := json.Marshal(result)
	if err != nil {
		return nil
	}
	_ = json.Unmarshal([]byte(marshal), v)
	return v
}

func (p *JsonRpcResult) ToBlock() *BlockResult {
	if p.checkErr() != nil {
		return nil
	}
	if (p).Result == nil {
		return nil
	}
	result := (p).Result.(map[string]interface{})
	if len(result) == 0 {
		return nil
	}
	v := &BlockResult{}
	marshal, err := json.Marshal(result)
	if err != nil {
		return nil
	}
	_ = json.Unmarshal([]byte(marshal), v)
	return v
}

func (p *JsonRpcResult) ToPaymentQueryInfo() *PaymentQueryInfo {
	if p.checkErr() != nil {
		return nil
	}
	if (p).Result == nil {
		return nil
	}
	result := (p).Result.(map[string]interface{})
	if len(result) == 0 {
		return nil
	}
	v := &PaymentQueryInfo{}
	marshal, err := json.Marshal(result)
	if err != nil {
		return nil
	}
	_ = json.Unmarshal([]byte(marshal), v)
	return v
}

func (p *JsonRpcResult) checkErr() error {
	if p.Error != nil {
		return errors.New(p.Error.Message)
	}
	return nil
}
