package model

import (
	"github.com/shopspring/decimal"
)

type OpenAccountJson struct {
	Address     string           `json:"address"`
	Balance     decimal.Decimal  `json:"balance"`
	KtonBalance *decimal.Decimal `json:"kton_balance,omitempty"`
	KtonLock    *decimal.Decimal `json:"kton_lock,omitempty"`
	Lock        decimal.Decimal  `json:"lock"`
}

type OpenExtrinsicsJson struct {
	BlockNum           int              `json:"block_num"`
	BlockTimestamp     int              `json:"block_timestamp"`
	CallModule         string           `json:"call_module"`
	CallModuleFunction string           `json:"call_module_function"`
	ExtrinsicIndex     string           `json:"extrinsic_index"`
	Finalized          bool             `json:"finalized"`
	From               string           `json:"from"`
	Hash               string           `json:"extrinsic_hash"`
	Params             []ExtrinsicParam `json:"params"`
	Success            bool             `json:"success"`
}

type OpenExtrinsicJson struct {
	BlockNum           int              `json:"block_num"`
	BlockTimestamp     int              `json:"block_timestamp"`
	CallModule         string           `json:"call_module"`
	CallModuleFunction string           `json:"call_module_function"`
	ExtrinsicIndex     string           `json:"extrinsic_index"`
	Event              *[]ChainEvent    `json:"event"`
	Finalized          bool             `json:"finalized"`
	From               string           `json:"from"`
	Hash               string           `json:"extrinsic_hash"`
	Nonce              int              `json:"nonce"`
	Params             []ExtrinsicParam `json:"params"`
	Signature          string           `json:"signature"`
	Success            bool             `json:"success"`
}

type OpenBlockJson struct {
	BlockNum       int    `json:"block_num"`
	BlockTimestamp int    `json:"block_timestamp"`
	Hash           string `json:"hash"`
	ParentHash     string `json:"parent_hash"`
	StateRoot      string `json:"state_root"`
	ExtrinsicsRoot string `json:"extrinsics_root"`
	SpecVersion    int    `json:"spec_version"`
	Validator      string `json:"validator"`
	Finalized      bool   `json:"finalized"`
}
