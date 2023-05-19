package model

import (
	"fmt"

	"github.com/itering/subscan/plugins/storage"
	"github.com/itering/subscan/util"
	"github.com/shopspring/decimal"
)

// SplitTableBlockNum
var SplitTableBlockNum = 1000000

type ChainBlock struct {
	ID              uint   `gorm:"primary_key" json:"id"`
	BlockNum        int    `gorm:"uniqueIndex" json:"block_num"`
	BlockTimestamp  int    `json:"block_timestamp"`
	Hash            string `gorm:"uniqueIndex;default: null;size:100" json:"hash"`
	ParentHash      string `gorm:"default: null;size:100" json:"parent_hash"`
	StateRoot       string `gorm:"default: null;size:100" json:"state_root"`
	ExtrinsicsRoot  string `gorm:"default: null;size:100" json:"extrinsics_root"`
	Logs            string `json:"logs" gorm:"type:text;"`
	Extrinsics      string `json:"extrinsics" gorm:"type:MEDIUMTEXT;"`
	EventCount      int    `json:"event_count"`
	ExtrinsicsCount int    `json:"extrinsics_count"`
	Event           string `json:"event" gorm:"type:MEDIUMTEXT;"`
	SpecVersion     int    `json:"spec_version"`
	Validator       string `json:"validator"`
	CodecError      bool   `json:"codec_error"`
	Finalized       bool   `json:"finalized"`
}

func (c ChainBlock) TableName() string {
	if c.BlockNum/SplitTableBlockNum == 0 {
		return "chain_blocks"
	}
	return fmt.Sprintf("chain_blocks_%d", c.BlockNum/SplitTableBlockNum)
}

func (c *ChainBlock) AsPlugin() *storage.Block {
	return &storage.Block{
		BlockNum:       c.BlockNum,
		BlockTimestamp: c.BlockTimestamp,
		Hash:           c.Hash,
		SpecVersion:    c.SpecVersion,
		Validator:      c.Validator,
		Finalized:      c.Finalized,
	}
}

type CallArg struct {
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

func (a CallArg) GetName() string {
	return a.Name
}

func (a CallArg) GetValue() interface{} {
	return a.Value
}

type ChainCall struct {
	BlockNum       int          `json:"block_num"`
	CallIdx        int          `json:"call_idx"`
	BlockTimestamp int          `json:"block_timestamp"`
	ExtrinsicHash  string       `json:"extrinsic_hash"`
	ModuleId       string       `json:"module_id"`
	CallId         string       `json:"call_id"`
	Params         []CallArg    `json:"params" sql:"type:text;"`
	Events         []ChainEvent `json:"events"`
}

type ChainEvent struct {
	ID            int         `gorm:"primary_key" json:"-"`
	EventIndex    string      `gorm:"index:event_idx,unique;index:event_index;default: null;size:100;" json:"event_index"`
	BlockNum      int         `gorm:"index" json:"block_num" `
	ExtrinsicIdx  int         `json:"extrinsic_idx"`
	Type          string      `gorm:"index" json:"-"`
	ModuleId      string      `gorm:"index" json:"module_id" `
	EventId       string      `gorm:"index" json:"event_id" `
	Params        interface{} `json:"params" gorm:"type:text;" `
	ExtrinsicHash string      `json:"extrinsic_hash" gorm:"default: null" `
	EventIdx      int         `gorm:"index:event_idx,unique" json:"event_idx"`
}

func (c ChainEvent) TableName() string {
	if c.BlockNum/SplitTableBlockNum == 0 {
		return "chain_events"
	}
	return fmt.Sprintf("chain_events_%d", c.BlockNum/SplitTableBlockNum)
}

type AsPlugin[P any] interface {
	AsPlugin() P
}

func MapAsPlugin[P interface{ *Pte }, T interface {
	AsPlugin[P]
	*Tp
}, Tp any, Pte any](arr []Tp,
) []Pte {
	return mapFunc(arr, func(item Tp) Pte {
		var t T = &item
		return *t.AsPlugin()
	})
}

func mapFunc[T any, R any](arr []T, f func(T) R) []R {
	mapped := make([]R, len(arr))
	for i, item := range arr {
		mapped[i] = f(item)
	}
	return mapped
}

func (c *ChainEvent) AsPlugin() *storage.Event {
	return &storage.Event{
		BlockNum:      c.BlockNum,
		ExtrinsicIdx:  c.ExtrinsicIdx,
		ModuleId:      c.ModuleId,
		EventId:       c.EventId,
		Params:        util.ToString(c.Params),
		ExtrinsicHash: c.ExtrinsicHash,
		EventIdx:      c.EventIdx,
		EventIndex:    c.EventIndex,
	}
}

func (c *CallArg) AsPlugin() *storage.CallArg {
	return &storage.CallArg{
		Name:  c.Name,
		Type:  c.Type,
		Value: c.Value,
	}
}

var _ AsPlugin[*storage.CallArg] = (*CallArg)(nil)

func (c *ChainCall) AsPlugin() *storage.Call {
	return &storage.Call{
		BlockNum:       c.BlockNum,
		CallIdx:        c.CallIdx,
		BlockTimestamp: c.BlockTimestamp,
		ExtrinsicHash:  c.ExtrinsicHash,
		ModuleId:       c.ModuleId,
		CallId:         c.CallId,
		Params:         MapAsPlugin[*storage.CallArg](c.Params),
		Events:         MapAsPlugin[*storage.Event](c.Events),
	}
}

type ChainExtrinsic struct {
	ID                 uint            `gorm:"primary_key"`
	ExtrinsicIndex     string          `gorm:"uniqueIndex;default: null;size:100" json:"extrinsic_index"`
	BlockNum           int             `gorm:"index" json:"block_num" `
	BlockTimestamp     int             `json:"block_timestamp"`
	ExtrinsicLength    string          `json:"extrinsic_length"`
	VersionInfo        string          `json:"version_info"`
	CallCode           string          `json:"call_code"`
	CallModuleFunction string          `gorm:"index;size:100" json:"call_module_function"`
	CallModule         string          `gorm:"index;size:100" json:"call_module"`
	Params             interface{}     `json:"params" gorm:"type:MEDIUMTEXT;" `
	AccountId          string          `gorm:"index:account_id" json:"account_id"`
	Signature          string          `json:"signature"`
	Nonce              int             `json:"nonce"`
	Era                string          `json:"era"`
	ExtrinsicHash      string          `gorm:"index;default: null" json:"extrinsic_hash"`
	IsSigned           bool            `gorm:"index:account_id;index" json:"is_signed"`
	Success            bool            `json:"success"`
	Fee                decimal.Decimal `json:"fee" gorm:"type:decimal(30,0);"`
	BatchIndex         int             `json:"-" gorm:"-"`
}

func (c ChainExtrinsic) TableName() string {
	if c.BlockNum/SplitTableBlockNum == 0 {
		return "chain_extrinsics"
	}
	return fmt.Sprintf("chain_extrinsics_%d", c.BlockNum/SplitTableBlockNum)
}

func (c *ChainExtrinsic) AsPlugin() *storage.Extrinsic {
	return &storage.Extrinsic{
		ExtrinsicIndex:     c.ExtrinsicIndex,
		CallModule:         c.CallModule,
		CallModuleFunction: c.CallModuleFunction,
		Params:             []byte(util.ToString(c.Params)),
		AccountId:          c.AccountId,
		Signature:          c.Signature,
		Nonce:              c.Nonce,
		Era:                c.Era,
		ExtrinsicHash:      c.ExtrinsicHash,
		Success:            c.Success,
		Fee:                c.Fee,
	}
}

type RuntimeVersion struct {
	Id          int    `json:"-"`
	Name        string `json:"-"`
	SpecVersion int    `gorm:"uniqueIndex" json:"spec_version"`
	Modules     string `json:"modules"  sql:"type:TEXT;"`
	RawData     string `json:"-" sql:"type:MEDIUMTEXT;"`
}

type RuntimeConstant struct {
	ID           uint   `gorm:"primary_key"`
	SpecVersion  int    `gorm:"index" json:"spec_version"`
	ModuleName   string `gorm:"index" json:"module_name" sql:"type:varchar(100);"`
	ConstantName string `gorm:"index" json:"constant_name" sql:"type:varchar(100);"`
	Type         string `json:"type" sql:"type:varchar(100);"`
	Value        string `json:"value" sql:"type:MEDIUMTEXT;"`
}

type ChainLog struct {
	ID        uint   `gorm:"primary_key"`
	BlockNum  int    `gorm:"index" json:"block_num" `
	LogIndex  string `gorm:"index;default: null;size:100" json:"log_index"`
	LogType   string `json:"log_type" `
	Data      string `json:"data" gorm:"type:text;"`
	Finalized bool   `json:"finalized"`
}

func (c ChainLog) TableName() string {
	if c.BlockNum/SplitTableBlockNum == 0 {
		return "chain_logs"
	}
	return fmt.Sprintf("chain_logs_%d", c.BlockNum/SplitTableBlockNum)
}

type ExtrinsicParam struct {
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	Value    interface{} `json:"value"`
	ValueRaw string      `json:"valueRaw"`
}

type EventParam struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}
