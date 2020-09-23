package model

type ExtrinsicError struct {
	ID            uint   `gorm:"primary_key" json:"-"`
	ExtrinsicHash string `json:"-" sql:"size:100;"`
	Module        string `json:"module"`
	Name          string `json:"name"`
	Doc           string `json:"doc"`
}

type MetadataModuleError struct {
	Module string   `json:"module"`
	Name   string   `json:"name"`
	Doc    []string `json:"doc"`
}
