package model

import (
	"fmt"
	"github.com/shopspring/decimal"
	"time"
)

// SplitTableBlockNum
var SplitTableBlockNum = 1000000

type ChainBlock struct {
	ID              uint      `gorm:"primary_key" json:"id"`
	BlockNum        int       `json:"block_num"`
	BlockTimestamp  int       `json:"block_timestamp"`
	CreatedAt       time.Time `json:"created_at"`
	Hash            string    `sql:"default: null;size:100" json:"hash"`
	ParentHash      string    `sql:"default: null;size:100" json:"parent_hash"`
	StateRoot       string    `sql:"default: null;size:100" json:"state_root"`
	ExtrinsicsRoot  string    `sql:"default: null;size:100" json:"extrinsics_root"`
	Logs            string    `json:"logs" sql:"type:text;"`
	Extrinsics      string    `json:"extrinsics" sql:"type:MEDIUMTEXT;"`
	EventCount      int       `json:"event_count"`
	ExtrinsicsCount int       `json:"extrinsics_count"`
	Event           string    `json:"event" sql:"type:text;"`
	SpecVersion     int       `json:"spec_version"`
	Validator       string    `json:"validator"`
	CodecError      bool      `json:"codec_error"`
	Finalized       bool      `json:"finalized"`
}

func (c ChainBlock) TableName() string {
	if c.BlockNum/SplitTableBlockNum == 0 {
		return "chain_blocks"
	}
	return fmt.Sprintf("chain_blocks_%d", c.BlockNum/SplitTableBlockNum)
}

type ChainEvent struct {
	ID            uint        `gorm:"primary_key" json:"-"`
	CreatedAt     time.Time   `json:"-" `
	EventIndex    string      `sql:"default: null;size:100;" json:"event_index"`
	BlockNum      int         `json:"block_num" `
	Phase         int         `json:"-"`
	ExtrinsicIdx  int         `json:"extrinsic_idx"`
	Type          string      `json:"-"`
	ModuleId      string      `json:"module_id" `
	EventId       string      `json:"event_id" `
	Params        interface{} `json:"params" sql:"type:text;" `
	ExtrinsicHash string      `json:"extrinsic_hash" sql:"default: null" `
	EventIdx      int         `json:"event_idx"`
	Finalized     bool        `json:"finalized"`
}

func (c ChainEvent) TableName() string {
	if c.BlockNum/SplitTableBlockNum == 0 {
		return "chain_events"
	}
	return fmt.Sprintf("chain_events_%d", c.BlockNum/SplitTableBlockNum)
}

type ChainExtrinsic struct {
	ID                 uint            `gorm:"primary_key"`
	CreatedAt          time.Time       `json:"created_at"`
	ExtrinsicIndex     string          `json:"extrinsic_index" sql:"default: null;size:100"`
	BlockNum           int             `json:"block_num" `
	BlockTimestamp     int             `json:"block_timestamp"`
	ValueRaw           string          `json:"value_raw"`
	ExtrinsicLength    string          `json:"extrinsic_length"`
	VersionInfo        string          `json:"version_info"`
	CallCode           string          `json:"call_code"`
	CallModuleFunction string          `json:"call_module_function"  sql:"size:100"`
	CallModule         string          `json:"call_module"  sql:"size:100"`
	Params             interface{}     `json:"params" sql:"type:MEDIUMTEXT;" `
	AccountLength      string          `json:"account_length"`
	AccountId          string          `json:"account_id"`
	AccountIndex       string          `json:"account_index"`
	Signature          string          `json:"signature"`
	Nonce              int             `json:"nonce"`
	Era                string          `json:"era"`
	ExtrinsicHash      string          `json:"extrinsic_hash" sql:"default: null" `
	IsSigned           bool            `json:"is_signed"`
	Success            bool            `json:"success"`
	Fee                decimal.Decimal `json:"fee" sql:"type:decimal(30,0);"`
	Finalized          bool            `json:"finalized"`
	BatchIndex         int             `json:"-" gorm:"-"`
}

func (c ChainExtrinsic) TableName() string {
	if c.BlockNum/SplitTableBlockNum == 0 {
		return "chain_extrinsics"
	}
	return fmt.Sprintf("chain_extrinsics_%d", c.BlockNum/SplitTableBlockNum)
}

type RuntimeVersion struct {
	Id          int    `json:"-"`
	Name        string `json:"-"`
	SpecVersion int    `json:"spec_version"`
	Modules     string `json:"modules"`
	RawData     string `json:"-" sql:"type:MEDIUMTEXT;"`
}

type ChainLog struct {
	ID         uint      `gorm:"primary_key"`
	CreatedAt  time.Time `json:"created_at"`
	BlockNum   int       `json:"block_num" `
	LogIndex   string    `json:"log_index" sql:"default: null;size:100"`
	LogType    string    `json:"log_type" `
	OriginType string    `json:"origin_type"`
	Data       string    `json:"data" sql:"type:text;"`
	Finalized  bool      `json:"finalized"`
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
	Type     string      `json:"type"`
	Value    interface{} `json:"value"`
	ValueRaw string      `json:"valueRaw"`
}
