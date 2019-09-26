package model

import "github.com/shopspring/decimal"

type ChainBlockJson struct {
	BlockNum        int                   `json:"block_num"`
	BlockTimestamp  int                   `json:"block_timestamp"`
	Hash            string                `json:"hash"`
	ParentHash      string                `json:"parent_hash"`
	StateRoot       string                `json:"state_root"`
	ExtrinsicsRoot  string                `json:"extrinsics_root"`
	Extrinsics      *[]ChainExtrinsicJson `json:"extrinsics"`
	Events          *[]ChainEventJson     `json:"events"`
	Logs            *[]ChainLogJson       `json:"logs"`
	EventCount      int                   `json:"event_count"`
	ExtrinsicsCount int                   `json:"extrinsics_count"`
	SpecVersion     int                   `json:"spec_version"`
	Validator       string                `json:"validator"`
	ValidatorName   string                `json:"validator_name"`
}

type SampleBlockJson struct {
	BlockNum        int    `json:"block_num"`
	BlockTimestamp  int    `json:"block_timestamp"`
	Hash            string `json:"hash"`
	EventCount      int    `json:"event_count"`
	ExtrinsicsCount int    `json:"extrinsics_count"`
	Validator       string `json:"validator"`
	ValidatorName   string `json:"validator_name"`
}

type ChainExtrinsicJson struct {
	BlockTimestamp     int    `json:"block_timestamp"`
	BlockNum           int    `json:"block_num"`
	ExtrinsicIndex     string `json:"extrinsic_index"`
	ValueRaw           string `json:"value_raw"`
	ExtrinsicLength    string `json:"extrinsic_length"`
	VersionInfo        string `json:"version_info"`
	CallCode           string `json:"call_code"`
	CallModuleFunction string `json:"call_module_function"`
	CallModule         string `json:"call_module"`
	Params             string `json:"params"`
	AccountLength      string `json:"account_length"`
	AccountId          string `json:"account_id"`
	AccountIndex       string `json:"account_index"`
	Signature          string `json:"signature"`
	Nonce              int    `json:"nonce"`
	Era                string `json:"era"`
	ExtrinsicHash      string `json:"extrinsic_hash"`
	Success            bool   `json:"success"`
}

type ExtrinsicDetail struct {
	BlockTimestamp     int               `json:"block_timestamp"`
	BlockNum           int               `json:"block_num"`
	ExtrinsicIndex     string            `json:"extrinsic_index"`
	CallModuleFunction string            `json:"call_module_function"`
	CallModule         string            `json:"call_module"`
	AccountId          string            `json:"account_id"`
	Signature          string            `json:"signature"`
	Nonce              int               `json:"nonce"`
	ExtrinsicHash      string            `json:"extrinsic_hash"`
	Success            bool              `json:"success"`
	Params             *[]ExtrinsicParam `json:"params"`
	Transfer           *TransferJson     `json:"transfer"`
	Event              *[]ChainEventJson `json:"event"`
}

type ChainEventJson struct {
	EventIndex    string `json:"event_index"`
	BlockNum      int    `json:"block_num"`
	Phase         int    `json:"phase"`
	ExtrinsicIdx  int    `json:"extrinsic_idx"`
	Type          string `json:"type"`
	ModuleId      string `json:"module_id"`
	EventId       string `json:"event_id"`
	Params        string `json:"params"`
	EventIdx      int    `json:"event_idx"`
	ExtrinsicHash string `json:"extrinsic_hash"`
}

type ChainLogJson struct {
	BlockNum   int    `json:"block_num" es:"type:keyword"`
	LogIndex   string `json:"log_index" sql:"default: null;size:100" es:"type:keyword"`
	LogType    string `json:"log_type" es:"type:keyword"`
	OriginType string `json:"origin_type"`
	Data       string `json:"data"`
}

type ChainTransactionJson struct {
	BlockTimestamp     int    `json:"block_timestamp"`
	ExtrinsicIndex     string `json:"extrinsic_index"`
	Signature          string `json:"signature"`
	FromHex            string `json:"from_hex"`
	Destination        string `json:"destination"`
	Hash               string `json:"hash"`
	BlockNum           int    `json:"block_num"`
	CallCode           string `json:"call_code"`
	CallModuleFunction string `json:"call_module_function"`
	CallModule         string `json:"call_module"`
	Params             string `json:"params"`
	Success            bool   `json:"success"`
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
}

type AccountJson struct {
	Address        string          `json:"address"`
	Balance        decimal.Decimal `json:"balance"`
	KtonBalance    decimal.Decimal `json:"kton_balance"`
	CountExtrinsic int             `json:"count_extrinsic"`
	Nonce          int             `json:"nonce"`
	AccountIndex   int             `json:"account_index"`
}
