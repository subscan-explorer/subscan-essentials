package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/itering/subscan-plugin/storage"
	"github.com/itering/subscan/util"
	"github.com/itering/subscan/util/address"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"strings"
)

const (
	IdGenerateCoefficient      = 100_000   // 100 thousand
	SplitTableBlockNum    uint = 1_000_000 // 1 million
)

type ChainBlock struct {
	ID              uint   `gorm:"primary_key" json:"id"`
	BlockNum        uint   `json:"block_num" gorm:"index:block_num,unique"`
	BlockTimestamp  int    `json:"block_timestamp"`
	Hash            string `gorm:"default: null;size:100;index:hash" json:"hash"`
	ParentHash      string `gorm:"default: null;size:100" json:"parent_hash"`
	StateRoot       string `gorm:"default: null;size:100" json:"state_root"`
	ExtrinsicsRoot  string `gorm:"default: null;size:100" json:"extrinsics_root"`
	EventCount      int    `json:"event_count"`
	ExtrinsicsCount int    `json:"extrinsics_count"`
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
		BlockNum:       int(c.BlockNum),
		BlockTimestamp: c.BlockTimestamp,
		Hash:           c.Hash,
		SpecVersion:    c.SpecVersion,
		Validator:      c.Validator,
		Finalized:      c.Finalized,
	}
}

type ChainEvent struct {
	ID             uint        `gorm:"primary_key;autoIncrement:false" json:"-"`
	ExtrinsicIndex string      `gorm:"default: null;size:100;" json:"extrinsic_index"`
	BlockNum       uint        `json:"block_num"  gorm:"index:block_num"`
	ExtrinsicIdx   int         `json:"extrinsic_idx"`
	ModuleId       string      `json:"module_id" gorm:"size:255;index:query_function"`
	EventId        string      `json:"event_id" gorm:"size:255;index:query_function"`
	Params         EventParams `json:"params" gorm:"type:json"`
	EventIdx       uint        `json:"event_idx"`
	Phase          int         `json:"phase" gorm:"size:8"`
}

func (c ChainEvent) TableName() string {
	if c.BlockNum/SplitTableBlockNum == 0 {
		return "chain_events"
	}
	return fmt.Sprintf("chain_events_%d", c.BlockNum/SplitTableBlockNum)
}

func (c ChainEvent) Id() uint {
	return c.BlockNum*IdGenerateCoefficient + c.EventIdx
}

func (c *ChainEvent) AsPlugin() *storage.Event {
	return &storage.Event{
		BlockNum:     int(c.BlockNum),
		ExtrinsicIdx: c.ExtrinsicIdx,
		ModuleId:     c.ModuleId,
		EventId:      c.EventId,
		Params:       c.Params.Marshal(),
		EventIdx:     int(c.EventIdx),
		Id:           c.Id(),
	}
}

type EventParams []EventParam

type EventParam struct {
	Type     string      `json:"type"`
	Value    interface{} `json:"value"`
	TypeName string      `json:"type_name,omitempty"`
}

func (j EventParams) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j EventParams) Marshal() []byte {
	b, _ := json.Marshal(j)
	return b
}

func (j *EventParams) Scan(src interface{}) error { return json.Unmarshal(src.([]byte), j) }

type ChainExtrinsic struct {
	ID                 uint   `gorm:"primary_key;autoIncrement:false"`
	ExtrinsicIndex     string `json:"extrinsic_index" gorm:"default: null;size:255;index:extrinsic_index"`
	BlockNum           uint   `json:"block_num" gorm:"index:block_num"`
	BlockTimestamp     int    `json:"block_timestamp"`
	CallModuleFunction string `json:"call_module_function"  gorm:"size:255;index:query_function"`
	CallModule         string `json:"call_module"  gorm:"size:255;index:query_function"`

	Params        ExtrinsicParams `json:"params" gorm:"type:json;"`
	AccountId     string          `json:"account_id"`
	Signature     string          `json:"signature"`
	Nonce         int             `json:"nonce"`
	Era           string          `json:"era"`
	ExtrinsicHash string          `json:"extrinsic_hash" gorm:"default: null;index:extrinsic_hash"`
	IsSigned      bool            `json:"is_signed"`
	Success       bool            `json:"success"`
	Fee           decimal.Decimal `json:"fee" gorm:"type:decimal(65,0);"`
}

type ExtrinsicParams []ExtrinsicParam

func (j ExtrinsicParams) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}

	return json.Marshal(j)
}

func (j *ExtrinsicParams) Scan(src interface{}) error { return json.Unmarshal(src.([]byte), j) }

func (j *ExtrinsicParams) Marshal() []byte {
	b, _ := json.Marshal(j)
	return b
}

type ExtrinsicParam struct {
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	Value    interface{} `json:"value"`
	TypeName string      `json:"type_name,omitempty"`
}

func ParsingExtrinsicParam(params interface{}) (param []ExtrinsicParam) {
	util.Logger().Error(util.UnmarshalAny(&param, params))
	return
}

func (c ChainExtrinsic) Id() uint {
	return ParseExtrinsicOrEventIndex(c.ExtrinsicIndex).GenerateId()
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
		Params:             c.Params.Marshal(),
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
	SpecVersion int    `json:"spec_version" gorm:"index:spec_version,unique"`
	Modules     string `json:"modules"  gorm:"type:TEXT;"`
	RawData     string `json:"-" gorm:"type:MEDIUMTEXT;"`
}

type ChainLog struct {
	ID        uint    `gorm:"primary_key"`
	BlockNum  uint    `json:"block_num" `
	LogIndex  string  `json:"log_index" gorm:"default: null;size:100"`
	LogType   string  `json:"log_type" `
	Data      LogData `json:"data" gorm:"type:json;"`
	Finalized bool    `json:"finalized"`
}

func (c ChainLog) TableName() string {
	if c.BlockNum/SplitTableBlockNum == 0 {
		return "chain_logs"
	}
	return fmt.Sprintf("chain_logs_%d", c.BlockNum/SplitTableBlockNum)
}

func (c ChainLog) Id() uint {
	return ParseExtrinsicOrEventIndex(c.LogIndex).GenerateId()
}

type LogData map[string]interface{}

func (l *LogData) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		if len(v) == 0 {
			return nil
		}
	case string:
		if v == "" {
			return nil
		}
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	var result LogData
	err := util.UnmarshalAny(&result, value)
	*l = result
	return err
}

func (l LogData) Value() (driver.Value, error) {
	return json.Marshal(l)
}

func (l LogData) Bytes() []byte {
	b, _ := json.Marshal(l)
	return b
}

type ExtrinsicOrEventIndex struct {
	BlockNum uint
	Index    uint
}

func ParseExtrinsicOrEventIndex(indexStr string) *ExtrinsicOrEventIndex {
	if sliceIndex := strings.Split(indexStr, "-"); len(sliceIndex) == 2 {
		return &ExtrinsicOrEventIndex{BlockNum: util.StringToUInt(sliceIndex[0]), Index: util.StringToUInt(sliceIndex[1])}
	}
	return nil
}

func ParseIndexInt(index uint) *ExtrinsicOrEventIndex {
	return &ExtrinsicOrEventIndex{BlockNum: index / IdGenerateCoefficient, Index: index % IdGenerateCoefficient}
}

func (e *ExtrinsicOrEventIndex) GenerateId() uint {
	if e == nil {
		return 0
	}
	return e.BlockNum*IdGenerateCoefficient + e.Index
}

func (e *ExtrinsicOrEventIndex) GenerateIndex() string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("%d-%d", e.BlockNum, e.Index)
}

func CheckoutParamValueAddress(value interface{}) string {
	switch a := value.(type) {
	case string:
		return address.Format(a)
	// Id         AccountId
	// Index      Compact<AccountIndex>
	// Raw   	  Bytes
	// Address32  H256
	// Address20  H160
	case map[string]interface{}: // multi address
		if v, ok := a["Id"]; ok {
			return address.Format(util.ToString(v))
		}
		if v, ok := a["Raw"]; ok {
			return address.Format(util.ToString(v))
		}
		if v, ok := a["Address32"]; ok {
			return address.Format(util.ToString(v))
		}
		if v, ok := a["Address20"]; ok {
			return address.Format(util.ToString(v))
		}
		if v, ok := a["Eth"]; ok {
			return address.Format(util.ToString(v))
		}
	}
	return address.Format(util.ToString(value))
}
