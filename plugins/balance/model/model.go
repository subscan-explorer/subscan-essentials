package model

import (
	"github.com/shopspring/decimal"
)

type Account struct {
	ID      uint            `gorm:"primary_key" json:"-"`
	Address string          `gorm:"default: null;size:100;index:address,unique" json:"address"`
	Nonce   int             `json:"nonce"`
	Balance decimal.Decimal `json:"balance" gorm:"type:decimal(65,0);"`
	Locked  decimal.Decimal `json:"locked" gorm:"type:decimal(65,0);"`
	Reserve decimal.Decimal `json:"reserve" gorm:"type:decimal(65,0);"`
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
