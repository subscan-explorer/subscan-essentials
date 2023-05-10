package dao

import (
	"github.com/itering/subscan/plugins/staking/model"
	"github.com/itering/subscan/plugins/storage"
)

func NewEraInfo(db storage.DB, info *model.EraInfo) error {
	return db.Create(info)
}

func GetEraInfoList(db storage.DB, page, row int) ([]model.EraInfo, int) {
	var list []model.EraInfo
	count, _ := db.FindBy(&list, nil, &storage.Option{
		PluginPrefix: "staking",
		Page:         page,
		PageSize:     row,
		Order:        "era DESC",
	})
	return list, count
}
