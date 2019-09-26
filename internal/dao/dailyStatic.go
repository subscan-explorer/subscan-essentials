package dao

import (
	"context"
	"github.com/jinzhu/gorm"
	"subscan-end/internal/model"
	"time"
)

func (d *Dao) IncrStatTransfer(c context.Context, utcTime time.Time) {
	var ds model.DailyStatic
	d.db.FirstOrCreate(&ds, &model.DailyStatic{TimeUTC: utcTime})
	d.db.Model(ds).Update(map[string]interface{}{"transfer_count": gorm.Expr("transfer_count + ?", 1)})
	_ = d.IncrMetadata(c, "count_transfer", 1)
	return
}

func (d *Dao) StatList(c context.Context, start, end string) *[]model.DailyStatic {
	var ds []model.DailyStatic
	d.db.Model(model.DailyStatic{}).Where("time_utc BETWEEN ? and ?", start, end).Find(&ds)
	return &ds
}
