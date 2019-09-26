package substrate

import (
	"encoding/json"
	"errors"
)

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

type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type SubParams struct {
	Result       interface{} `json:"result"`
	Subscription int64       `json:"subscription,omitempty"`
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
		return "", errors.New("no results")
	}
	return p.Result.(string), nil
}

func (p *JsonRpcResult) ToInt() uint64 {
	if p.checkErr() != nil {
		return 0
	}
	return p.Result.(uint64)
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

func (p *JsonRpcResult) ToToken() *SystemTokenResult {
	if p.checkErr() != nil {
		return nil
	}
	result := (p).Result.(map[string]interface{})
	if len(result) == 0 {
		return nil
	}
	v := &SystemTokenResult{}
	marshal, err := json.Marshal(result)
	if err != nil {
		return nil
	}
	_ = json.Unmarshal([]byte(marshal), v)
	return v
}

func (p *JsonRpcResult) ToStorage() *StateStorageResult {
	if p.Params == nil {
		return nil
	}
	result := (p).Params.Result.(map[string]interface{})
	if len(result) == 0 {
		return nil
	}
	v := &StateStorageResult{}
	marshal, err := json.Marshal(result)
	if err != nil {
		return nil
	}
	_ = json.Unmarshal([]byte(marshal), v)
	return v
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

func (p *JsonRpcResult) checkErr() error {
	if p.Error != nil {
		return errors.New(p.Error.Message)
	}
	return nil
}

func (p *JsonRpcParams) checkErr() error {
	if p.Error != nil {
		return errors.New(p.Error.Message)
	}
	if p.Params == nil {
		return errors.New("nil result")
	}
	return nil
}
