package dao

import (
	"github.com/itering/subscan/model"
	"github.com/pkg/errors"
)

func (d *Dao) CreateRuntimeConstants(spec int, constants []model.RuntimeConstant) error {
	for _, constant := range constants {
		q := d.db.Save(&constant)
		if q.RowsAffected == 0 {
			return errors.New("create runtime constant failed")
		}
	}
	return nil
}

func (d *ReadOnlyDao) GetRuntimeConstantLatest(moduleName string, constantName string) *model.RuntimeConstant {
	constant, _ := findOne[model.RuntimeConstant](d, "*", where("module_name = ? AND constant_name = ?", moduleName, constantName), "spec_version DESC")
	return constant
}
