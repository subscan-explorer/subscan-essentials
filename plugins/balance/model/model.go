package model

import (
	"github.com/shopspring/decimal"
)

type Account struct {
	ID       uint            `gorm:"primary_key" json:"-"`
	Address  string          `gorm:"default: null;size:100;index:address,unique" json:"address"`
	Nonce    int             `json:"nonce"`
	Balance  decimal.Decimal `json:"balance" gorm:"type:decimal(65,0);index:balance"`
	Locked   decimal.Decimal `json:"locked" gorm:"type:decimal(65,0);"`
	Reserved decimal.Decimal `json:"reserved" gorm:"type:decimal(65,0);"`
}

func (a *Account) TableName() string {
	return "balance_accounts"
}

type AccountData struct {
	Nonce    int `json:"nonce"`
	RefCount int `json:"ref_count"`
	Data     struct {
		Free       decimal.Decimal `json:"free"`
		Reserved   decimal.Decimal `json:"reserved"`
		MiscFrozen decimal.Decimal `json:"miscFrozen"`
		FeeFrozen  decimal.Decimal `json:"feeFrozen"`
	} `json:"data"`
}

type Transfer struct {
	Id             uint            `json:"id" gorm:"primary_key;autoIncrement:false"`
	BlockNum       uint            `json:"blockNum" gorm:"size:32"`
	Sender         string          `json:"sender" gorm:"size:255;index:query_function"`
	Receiver       string          `json:"receiver" gorm:"size:255;index:query_function"`
	Amount         decimal.Decimal `json:"amount" gorm:"decimal(65)"`
	BlockTimestamp int64           `json:"block_timestamp" `
	Symbol         string          `json:"symbol" gorm:"size:255"`
	TokenId        string          `json:"token_id" gorm:"size:255"`
	ExtrinsicIndex string          `json:"extrinsic_index" gorm:"size:255;index:extrinsic_index"`
}

func (a *Transfer) TableName() string {
	return "balance_transfers"
}
