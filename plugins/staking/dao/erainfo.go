package dao

import (
	"github.com/itering/subscan/plugins/staking/model"
	"github.com/itering/subscan/plugins/storage"
	"github.com/itering/subscan/util/address"
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

type EraAndPoints struct {
	Era    uint32
	Points map[address.SS58Address]uint32
}

func GetEraPointsList(db storage.DB, page, row int) ([]EraAndPoints, int) {
	var list []model.EraInfo
	db.Query(&model.EraInfo{}).Select("era, validator_points").Order("era DESC").Limit(row).Offset((page - 1) * row).Find(&list)
	var result []EraAndPoints

	for _, item := range list {
		result = append(result, EraAndPoints{Era: item.Era, Points: item.ValidatorPoints.Data()})
	}

	return result, len(result)
}
