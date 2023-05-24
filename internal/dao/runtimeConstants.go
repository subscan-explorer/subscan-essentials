package dao

import (
	"github.com/itering/subscan/model"
	"github.com/pkg/errors"
	"golang.org/x/exp/slog"
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
	constant, err := findOne[model.RuntimeConstant](d, "*", where("module_name = ? AND constant_name = ?", moduleName, constantName), "spec_version DESC")
	if err != nil {
		slog.Error("get runtime constant latest failed", "error", err)
		return nil
	}
	return constant
}
