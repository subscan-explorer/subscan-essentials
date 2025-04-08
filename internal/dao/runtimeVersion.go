package dao

import (
	"github.com/itering/subscan/model"
	"github.com/itering/substrate-api-rpc/metadata"
)

func (d *Dao) CreateRuntimeVersion(name string, specVersion int) int64 {
	query := d.db.Scopes(model.IgnoreDuplicate).Create(&model.RuntimeVersion{
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
	var list model.RuntimeVersion
	query := d.db.Select("spec_version,raw_data").Model(model.RuntimeVersion{}).Order("spec_version DESC").First(&list)
	if query.Error != nil {
		return nil
	}
	return &list
}

func (d *Dao) RuntimeVersionRaw(spec int) *metadata.RuntimeRaw {
	var one model.RuntimeVersion
	query := d.db.Model(model.RuntimeVersion{}).Where("spec_version = ?", spec).First(&one)
	if query.Error != nil {
		return nil
	}
	return &metadata.RuntimeRaw{
		Spec: one.SpecVersion,
		Raw:  one.RawData,
	}
}
