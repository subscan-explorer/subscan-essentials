package dao

import (
	"context"
	"subscan-end/internal/model"
)

func (d *Dao) CreateRuntimeVersion(c context.Context, name string, specVersion int) error {
	query := d.db.Create(&model.RuntimeVersion{
		Id:          specVersion,
		Name:        name,
		SpecVersion: specVersion,
	})
	return query.Error
}
