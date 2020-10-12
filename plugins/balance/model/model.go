package model

import (
	"github.com/shopspring/decimal"
)

type Account struct {
	ID      uint            `gorm:"primary_key" json:"-"`
	Address string          `sql:"default: null;size:100" json:"address"`
	Nonce   int             `json:"nonce"`
	Balance decimal.Decimal `json:"balance" sql:"type:decimal(30,0);"`
	Lock    decimal.Decimal `json:"lock" sql:"type:decimal(30,0);"`
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
