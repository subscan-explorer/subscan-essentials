package dao

import (
	"github.com/itering/subscan/model"
	"github.com/pkg/errors"
)

func (d *Dao) CreateRuntimeConstants(spec int, constants []model.RuntimeConstant) error {
	for _, constant := range constants {
		q := d.db.Create(&constant)
		if q.RowsAffected == 0 {
			return errors.New("create runtime constant failed")
		}
	}
	return nil
}

func (d *Dao) GetRuntimeConstantLatest(moduleName string, constantName string) *model.RuntimeConstant {
	var constant model.RuntimeConstant
	d.db.Where("module_name = ? AND constant_name = ?", moduleName, constantName).Order("spec_version DESC").First(&constant)
	if constant == (model.RuntimeConstant{}) {
		return nil
	}
	return &constant
}
