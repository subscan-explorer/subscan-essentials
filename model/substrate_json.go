package model

import "github.com/shopspring/decimal"

type ChainBlockJson struct {
	BlockNum          int                  `json:"block_num"`
	BlockTimestamp    int                  `json:"block_timestamp"`
	Hash              string               `json:"hash"`
	ParentHash        string               `json:"parent_hash"`
	StateRoot         string               `json:"state_root"`
	ExtrinsicsRoot    string               `json:"extrinsics_root"`
	Extrinsics        []ChainExtrinsicJson `json:"extrinsics"`
	Events            []ChainEventJson     `json:"events"`
	Logs              []ChainLogJson       `json:"logs"`
	EventCount        int                  `json:"event_count"`
	ExtrinsicsCount   int                  `json:"extrinsics_count"`
	SpecVersion       int                  `json:"spec_version"`
	Validator         string               `json:"validator"`
	ValidatorName     string               `json:"validator_name"`
	ValidatorIndexIds string               `json:"validator_index_ids"`
	Finalized         bool                 `json:"finalized"`
}

type SampleBlockJson struct {
	BlockNum          int    `json:"block_num"`
	BlockTimestamp    int    `json:"block_timestamp"`
	Hash              string `json:"hash"`
	EventCount        int    `json:"event_count"`
	ExtrinsicsCount   int    `json:"extrinsics_count"`
	Validator         string `json:"validator"`
	ValidatorName     string `json:"validator_name"`
	ValidatorIndexIds string `json:"validator_index_ids"`
	Finalized         bool   `json:"finalized"`
}

type ChainExtrinsicJson struct {
	BlockTimestamp     int             `json:"block_timestamp"`
	BlockNum           int             `json:"block_num"`
	ExtrinsicIndex     string          `json:"extrinsic_index"`
	CallModuleFunction string          `json:"call_module_function"`
	CallModule         string          `json:"call_module"`
	Params             string          `json:"params"`
	AccountId          string          `json:"account_id"`
	AccountIndex       string          `json:"account_index"`
	Signature          string          `json:"signature"`
	Nonce              int             `json:"nonce"`
	ExtrinsicHash      string          `json:"extrinsic_hash"`
	Success            bool            `json:"success"`
	Fee                decimal.Decimal `json:"fee"`
}

type ExtrinsicDetail struct {
	BlockTimestamp     int              `json:"block_timestamp"`
	BlockNum           int              `json:"block_num"`
	ExtrinsicIndex     string           `json:"extrinsic_index"`
	CallModuleFunction string           `json:"call_module_function"`
	CallModule         string           `json:"call_module"`
	AccountId          string           `json:"account_id"`
	Signature          string           `json:"signature"`
	Nonce              int              `json:"nonce"`
	ExtrinsicHash      string           `json:"extrinsic_hash"`
	Success            bool             `json:"success"`
	Params             []ExtrinsicParam `json:"params"`
	Event              *[]ChainEvent    `json:"event"`
	Fee                decimal.Decimal  `json:"fee"`
	Finalized          bool             `json:"finalized"`
}

type ChainEventJson struct {
	EventIndex     string `json:"event_index"`
	BlockNum       int    `json:"block_num"`
	ExtrinsicIdx   int    `json:"extrinsic_idx"`
	ModuleId       string `json:"module_id"`
	EventId        string `json:"event_id"`
	Params         string `json:"params"`
	EventIdx       int    `json:"event_idx"`
	ExtrinsicHash  string `json:"extrinsic_hash"`
	BlockTimestamp int    `json:"block_timestamp"`
}

type ChainLogJson struct {
	BlockNum   int    `json:"block_num" es:"type:keyword"`
	LogIndex   string `json:"log_index" sql:"default: null;size:100" es:"type:keyword"`
	LogType    string `json:"log_type" es:"type:keyword"`
	OriginType string `json:"origin_type"`
	Data       string `json:"data"`
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
