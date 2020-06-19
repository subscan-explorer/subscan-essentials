package model

import (
	"fmt"
	"github.com/shopspring/decimal"
	"time"
)

// SplitTableBlockNum
var SplitTableBlockNum = 1000000

type ChainMetadata struct {
	Version int    `json:"version"`
	Type    string `json:"type"` // storage/call/event/constant/error
	Name    string `json:"name"`
	Prefix  string `json:"prefix"`
}

type ChainBlock struct {
	ID              uint      `gorm:"primary_key" json:"id"`
	BlockNum        int       `json:"block_num" es:"type:keyword"`
	BlockTimestamp  int       `json:"block_timestamp"`
	CreatedAt       time.Time `json:"created_at" es:"type:date"`
	Hash            string    `sql:"default: null;size:100" json:"hash" es:"type:keyword"`
	ParentHash      string    `sql:"default: null;size:100" json:"parent_hash" es:"type:keyword"`
	StateRoot       string    `sql:"default: null;size:100" json:"state_root" es:"type:keyword"`
	ExtrinsicsRoot  string    `sql:"default: null;size:100" json:"extrinsics_root" es:"type:keyword"`
	Logs            string    `json:"logs" sql:"type:text;"`
	Extrinsics      string    `json:"extrinsics" sql:"type:MEDIUMTEXT;"`
	EventCount      int       `json:"event_count"`
	ExtrinsicsCount int       `json:"extrinsics_count"`
	Event           string    `json:"event" sql:"type:text;"`
	SpecVersion     int       `json:"spec_version"  es:"type:text"`
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
	CreatedAt     time.Time   `json:"-" es:"type:date"`
	EventIndex    string      `sql:"default: null;size:100;" json:"event_index" es:"type:keyword"`
	BlockNum      int         `json:"block_num" es:"type:keyword"`
	Phase         int         `json:"-"`
	ExtrinsicIdx  int         `json:"extrinsic_idx"`
	Type          string      `json:"-"`
	ModuleId      string      `json:"module_id" es:"type:keyword"`
	EventId       string      `json:"event_id" es:"type:keyword"`
	Params        interface{} `json:"params" sql:"type:text;" es:"type:text"`
	ExtrinsicHash string      `json:"extrinsic_hash" sql:"default: null" es:"type:keyword"`
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
	ExtrinsicIndex     string          `json:"extrinsic_index" sql:"default: null;size:100" es:"type:keyword"`
	BlockNum           int             `json:"block_num" es:"type:keyword"`
	BlockTimestamp     int             `json:"block_timestamp"`
	ValueRaw           string          `json:"value_raw"`
	ExtrinsicLength    string          `json:"extrinsic_length"`
	VersionInfo        string          `json:"version_info"`
	CallCode           string          `json:"call_code"`
	CallModuleFunction string          `json:"call_module_function" es:"type:keyword" sql:"size:100"`
	CallModule         string          `json:"call_module" es:"type:keyword" sql:"size:100"`
	Params             interface{}     `json:"params" sql:"type:MEDIUMTEXT;" es:"type:text"`
	AccountLength      string          `json:"account_length"`
	AccountId          string          `json:"account_id"`
	AccountIndex       string          `json:"account_index"`
	Signature          string          `json:"signature"`
	Nonce              int             `json:"nonce"`
	Era                string          `json:"era"`
	ExtrinsicHash      string          `json:"extrinsic_hash" sql:"default: null" es:"type:keyword"`
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

type ChainTransaction struct {
	ID                 uint            `gorm:"primary_key"`
	CreatedAt          time.Time       `json:"created_at"`
	FromHex            string          `json:"from_hex" es:"type:keyword"`
	Destination        string          `json:"destination" es:"type:keyword"`
	ExtrinsicIndex     string          `sql:"default: null;size:100" json:"extrinsic_index" es:"type:keyword"`
	Signature          string          `json:"signature"`
	Success            bool            `json:"success"`
	Hash               string          `sql:"default: null;size:100" json:"hash" es:"type:keyword"`
	BlockNum           int             `json:"block_num" es:"type:keyword"`
	BlockTimestamp     int             `json:"block_timestamp"`
	CallCode           string          `json:"call_code"`
	CallModuleFunction string          `json:"call_module_function" es:"type:keyword" sql:"size:100"`
	CallModule         string          `json:"call_module" es:"type:keyword" sql:"size:100"`
	Params             interface{}     `json:"params" sql:"type:MEDIUMTEXT;" es:"type:text"`
	Amount             decimal.Decimal `json:"amount" sql:"type:decimal(30,15);"`
	Fee                decimal.Decimal `json:"fee" sql:"type:decimal(30,0);"`
	Finalized          bool            `json:"finalized"`
}

func (c ChainTransaction) TableName() string {
	if c.BlockNum/SplitTableBlockNum == 0 {
		return "chain_transactions"
	}
	return fmt.Sprintf("chain_transactions_%d", c.BlockNum/SplitTableBlockNum)
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

type RuntimeVersion struct {
	Id                    int    `json:"-"`
	Name                  string `json:"-"`
	SpecVersion           int    `json:"spec_version"`
	CountCallFunctions    int    `json:"-"`
	CountEvents           int    `json:"-"`
	CountStorageFunctions int    `json:"-"`
	CountConstants        int    `json:"-"`
	CountErrorType        int    `json:"-"`
	Modules               string `json:"modules"`
	RawData               string `json:"-" sql:"type:MEDIUMTEXT;"`
}

type ChainLog struct {
	ID         uint      `gorm:"primary_key"`
	CreatedAt  time.Time `json:"created_at"`
	BlockNum   int       `json:"block_num" es:"type:keyword"`
	LogIndex   string    `json:"log_index" sql:"default: null;size:100" es:"type:keyword"`
	LogType    string    `json:"log_type" es:"type:keyword"`
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

type ExtrinsicError struct {
	ID            uint   `gorm:"primary_key" json:"-"`
	ExtrinsicHash string `json:"-" sql:"size:100;"`
	Module        string `json:"module"`
	Name          string `json:"name"`
	Doc           string `json:"doc"`
}

type DispatchError struct {
	Other        *string     `json:"Other,omitempty"`
	CannotLookup *string     `json:"CannotLookup,omitempty"`
	BadOrigin    *string     `json:"BadOrigin,omitempty"`
	Module       interface{} `json:"Module,omitempty"`
	Error        *int        `json:"error,omitempty"`
}

type DispatchErrorModule struct {
	Index int `json:"index"`
	Error int `json:"error"`
}
