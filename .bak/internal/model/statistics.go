package model

import "time"

type DailyStatic struct {
	ID            uint      `gorm:"primary_key"`
	TimeUTC       time.Time `json:"time_utc" sql:"type:date;"`
	TransferCount int       `json:"transfer_count"`
}
