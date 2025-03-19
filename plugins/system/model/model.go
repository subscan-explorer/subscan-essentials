package model

type ExtrinsicError struct {
	ID             uint   `gorm:"primary_key" json:"-"`
	ExtrinsicIndex string `json:"-" gorm:"size:100;index:extrinsic_index,unique"`
	Module         string `json:"module"`
	Name           string `json:"name"`
	Doc            string `json:"doc"`
}

type MetadataModuleError struct {
	Module string   `json:"module"`
	Name   string   `json:"name"`
	Doc    []string `json:"doc"`
}
