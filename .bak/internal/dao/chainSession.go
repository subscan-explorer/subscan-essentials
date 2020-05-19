package dao

import (
	"context"
	"subscan-end/internal/model"
)

func (d *Dao) AddSessionValidator(c context.Context, cs *model.SessionValidator) error {
	query := d.db.Create(&cs)
	return query.Error
}

func (d *Dao) AddSessionNominator(c context.Context, sn *model.SessionNominator) error {
	query := d.db.Create(&sn)
	return query.Error
}

func (d *Dao) AddSession(c context.Context, cs *model.ChainSession) error {
	query := d.db.Create(&cs)
	return query.Error
}
