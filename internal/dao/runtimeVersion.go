package dao

import (
	"github.com/itering/subscan/model"
	"github.com/itering/substrate-api-rpc/metadata"
)

func (d *Dao) CreateRuntimeVersion(name string, specVersion int) int64 {
	var versions []model.RuntimeVersion
	d.db.Select("id").Where("spec_version = ?", specVersion).Limit(1).Find(&versions)
	if len(versions) > 0 {
		return 0
	}

	query := d.db.Create(&model.RuntimeVersion{
		Name:        name,
		SpecVersion: specVersion,
	})
	return query.RowsAffected
}

func (d *Dao) SetRuntimeData(specVersion int, modules string, rawData string) int64 {
	query := d.db.Model(model.RuntimeVersion{}).Where("spec_version=?", specVersion).Updates(model.RuntimeVersion{
		Modules: modules,
		RawData: rawData,
	})
	return query.RowsAffected
}

func (d *Dao) RuntimeVersionList() []model.RuntimeVersion {
	var list []model.RuntimeVersion
	d.db.Select("spec_version,modules").Model(model.RuntimeVersion{}).Find(&list)
	return list
}

func (d *Dao) RuntimeVersionRecent() *model.RuntimeVersion {
	var list []model.RuntimeVersion
	_ = d.db.Select("spec_version,raw_data").Model(model.RuntimeVersion{}).Limit(1).Order("spec_version DESC").Find(&list)
	if len(list) == 0 {
		return nil
	}
	return &list[0]
}

func (d *Dao) RuntimeVersionRaw(spec int) *metadata.RuntimeRaw {
	var one metadata.RuntimeRaw
	query := d.db.Model(model.RuntimeVersion{}).
		Select("spec_version as spec ,raw_data as raw").
		Where("spec_version = ?", spec).
		Scan(&one)
	if RecordNotFound(query) {
		return nil
	}
	return &one
}
