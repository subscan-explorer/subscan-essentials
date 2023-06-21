package model

import "gorm.io/gorm"

type LastProcessedBlock struct {
	gorm.Model
	Number int `gorm:"unique_index;not null"`
}
