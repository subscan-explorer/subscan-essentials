package model

import "github.com/shopspring/decimal"

type ChainBlockJson struct {
	BlockNum        uint   `json:"block_num"`
	BlockTimestamp  int    `json:"block_timestamp"`
	Hash            string `json:"hash"`
	ParentHash      string `json:"parent_hash"`
	StateRoot       string `json:"state_root"`
	ExtrinsicsRoot  string `json:"extrinsics_root"`
	EventCount      int    `json:"event_count"`
	ExtrinsicsCount int    `json:"extrinsics_count"`
	SpecVersion     int    `json:"spec_version"`
	Validator       string `json:"validator"`
	Finalized       bool   `json:"finalized"`
}

type SampleBlockJson struct {
	BlockNum        uint   `json:"block_num"`
	BlockTimestamp  int    `json:"block_timestamp"`
	Hash            string `json:"hash"`
	EventCount      int    `json:"event_count"`
	ExtrinsicsCount int    `json:"extrinsics_count"`
	Validator       string `json:"validator"`
	Finalized       bool   `json:"finalized"`
}

type ChainExtrinsicJson struct {
	Id                 uint            `json:"id"`
	BlockTimestamp     int             `json:"block_timestamp"`
	BlockNum           uint            `json:"block_num"`
	ExtrinsicIndex     string          `json:"extrinsic_index"`
	CallModuleFunction string          `json:"call_module_function"`
	CallModule         string          `json:"call_module"`
	Params             ExtrinsicParams `json:"params,omitempty"`
	AccountId          string          `json:"account_id"`
	Signature          string          `json:"signature"`
	Nonce              int             `json:"nonce"`
	ExtrinsicHash      string          `json:"extrinsic_hash"`
	Success            bool            `json:"success"`
	Fee                decimal.Decimal `json:"fee"`
}

type ExtrinsicDetail struct {
	BlockTimestamp     int             `json:"block_timestamp"`
	BlockNum           uint            `json:"block_num"`
	ExtrinsicIndex     string          `json:"extrinsic_index"`
	CallModuleFunction string          `json:"call_module_function"`
	CallModule         string          `json:"call_module"`
	AccountId          string          `json:"account_id"`
	Signature          string          `json:"signature"`
	Nonce              int             `json:"nonce"`
	ExtrinsicHash      string          `json:"extrinsic_hash"`
	Success            bool            `json:"success"`
	Params             ExtrinsicParams `json:"params"`
	Fee                decimal.Decimal `json:"fee"`
	Finalized          bool            `json:"finalized"`
	Lifetime           *Lifetime       `json:"lifetime"`
}

type Lifetime struct {
	Birth uint64 `json:"birth"`
	Death uint64 `json:"death"`
}

type ChainEventJson struct {
	Id             uint        `json:"id"`
	EventIndex     string      `json:"event_index"`
	ExtrinsicIndex string      `json:"extrinsic_index"`
	BlockNum       uint        `json:"block_num"`
	ModuleId       string      `json:"module_id"`
	EventId        string      `json:"event_id"`
	Params         EventParams `json:"params"`
	EventIdx       uint        `json:"event_idx"`
	BlockTimestamp int         `json:"block_timestamp"`
	Phase          int         `json:"phase"`
}

type ChainLogJson struct {
	BlockNum int    `json:"block_num"`
	LogIndex string `json:"log_index" gorm:"default: null;size:100"`
	LogType  string `json:"log_type"`
	Data     string `json:"data"`
}

type TransferJson struct {
	From           string          `json:"from"`
	To             string          `json:"to"`
	Module         string          `json:"module"`
	Amount         decimal.Decimal `json:"amount"`
	Hash           string          `json:"hash"`
	BlockTimestamp int             `json:"block_timestamp"`
	BlockNum       int             `json:"block_num"`
	ExtrinsicIndex string          `json:"extrinsic_index"`
	Success        bool            `json:"success"`
	Fee            decimal.Decimal `json:"fee"`
}

type ExtrinsicsJson struct {
	FromHex            string          `json:"from_hex"`
	BlockNum           int             `json:"block_num"`
	BlockTimestamp     int             `json:"block_timestamp"`
	ExtrinsicIndex     string          `json:"extrinsic_index"`
	Hash               string          `json:"extrinsic_hash"`
	Success            bool            `json:"success"`
	CallModule         string          `json:"call_module"`
	CallModuleFunction string          `json:"call_module_function"`
	Params             string          `json:"params"`
	Fee                decimal.Decimal `json:"fee,omitempty"`
	Destination        string          `json:"destination,omitempty"`
	Amount             decimal.Decimal `json:"amount,omitempty"`
	Finalized          bool            `json:"finalized"`
}

type EventRecord struct {
	Phase        int          `json:"phase"`
	ExtrinsicIdx int          `json:"extrinsic_idx"`
	Type         string       `json:"type"`
	ModuleId     string       `json:"module_id"`
	EventId      string       `json:"event_id"`
	Params       []EventParam `json:"params"`
	Topics       []string     `json:"topics"`
	EventIdx     int          `json:"event_idx"`
}
