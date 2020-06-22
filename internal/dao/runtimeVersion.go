package dao

import (
	"github.com/itering/subscan/internal/model"
	"github.com/itering/subscan/internal/substrate/metadata"
)

func (d *Dao) CreateRuntimeVersion(name string, specVersion int) int64 {
	query := d.Db.Create(&model.RuntimeVersion{
		Name:        name,
		SpecVersion: specVersion,
	})
	return query.RowsAffected
}

func (d *Dao) SetRuntimeData(specVersion int, modules string, rawData string) int64 {
	query := d.Db.Model(model.RuntimeVersion{}).Where("spec_version=?", specVersion).UpdateColumn(model.RuntimeVersion{
		Modules: modules,
		RawData: rawData,
	})
	return query.RowsAffected
}

func (d *Dao) RuntimeVersionList() []model.RuntimeVersion {
	var list []model.RuntimeVersion
	d.Db.Select("spec_version,modules").Model(model.RuntimeVersion{}).Find(&list)
	return list
}

func (d *Dao) RuntimeVersionRecent() *model.RuntimeVersion {
	var list model.RuntimeVersion
	query := d.Db.Select("spec_version,raw_data").Model(model.RuntimeVersion{}).Order("spec_version DESC").First(&list)
	if query.RecordNotFound() {
		return nil
	}
	return &list
}

func (d *Dao) RuntimeVersionRaws(spec int) *[]metadata.RuntimeRaw {
	var list []metadata.RuntimeRaw
	d.Db.Model(model.RuntimeVersion{}).
		Select("spec_version as spec ,raw_data as raw").
		Where("spec_version = ?", spec).
		Scan(&list)
	return &list
}
