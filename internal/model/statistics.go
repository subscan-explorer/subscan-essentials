package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type DailyStatic struct {
	ID             uint            `gorm:"primary_key"`
	TimeUTC        string          `json:"time_utc" sql:"type:date;"`
	TimeHourUTC    time.Time       `json:"time_hour_utc" sql:"type:datetime;"`
	TimeSixHourUTC time.Time       `json:"time_six_hour_utc" sql:"type:datetime;"`
	TransferCount  int             `json:"transfer_count"`
	ExtrinsicCount int             `json:"extrinsic_count"`
	TransferAmount decimal.Decimal `json:"transfer_amount"  sql:"type:decimal(60,20);"`
}

type DailyStaticJson struct {
	TimeUTC             string          `json:"time_utc"`
	TimeHourUTC         time.Time       `json:"time_hour_utc"`
	TimeSixHourUTC      time.Time       `json:"time_six_hour_utc"`
	Total               int             `json:"total"`
	TransferAmountTotal decimal.Decimal `json:"transfer_amount_total"`
}

type BlockOutputStat struct {
	Validator string `json:"validator"`
	Total     int    `json:"total"`
}

type EraBondStat struct {
	Era          int             `json:"era"`
	Owner        decimal.Decimal `json:"owner"`
	TotalBond    decimal.Decimal `json:"total_bond"`
	Avg          decimal.Decimal `json:"avg"`
	TotalBondAvg decimal.Decimal `json:"total_avg"`
}
