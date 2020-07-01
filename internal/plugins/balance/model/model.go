package model

import (
	"github.com/shopspring/decimal"
)

type Account struct {
	ID      uint            `gorm:"primary_key"`
	Address string          `sql:"default: null;size:100" json:"address"`
	Nonce   int             `json:"nonce"`
	Balance decimal.Decimal `json:"balance" sql:"type:decimal(30,0);"`
	Lock    decimal.Decimal `json:"lock" sql:"type:decimal(30,0);"`
}
