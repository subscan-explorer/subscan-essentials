package dao

import (
	"context"
	"fmt"
	"subscan-end/internal/model"
)

func (d *Dao) UpdateValidatorName(c context.Context, nodeName, controller string) error {
	v := model.ValidatorInfo{
		NodeName: nodeName,
	}
	query := d.db.Model(model.ValidatorInfo{}).Where("validator_controller = ?", controller).Where("node_name=?", "").Update(&v)
	return query.Error
}

// controller bind session
func (d *Dao) CreateValidatorSessionOrUpdate(c context.Context, controller, session string) error {
	if session == "" {
		return nil
	}
	var validator model.ValidatorInfo
	query := d.db.Model(model.ValidatorInfo{}).Where("validator_controller = ?", controller).First(&validator)
	if query.RecordNotFound() {
		query = d.db.Create(&model.ValidatorInfo{
			ValidatorController: controller,
			ValidatorSession:    session,
		})
	} else {
		query = d.db.Model(&validator).Update(&model.ValidatorInfo{ValidatorSession: session})
	}
	return query.Error
}

// stash bind controller
func (d *Dao) CreateValidatorStashOrUpdate(c context.Context, controller, stash string) error {
	if stash == "" {
		return nil
	}
	var validatorStash model.ValidatorInfo
	var validatorController model.ValidatorInfo
	query := d.db.Model(model.ValidatorInfo{}).Where("validator_stash = ?", stash).First(&validatorStash)
	queryController := d.db.Model(model.ValidatorInfo{}).Where("validator_controller = ?", controller).First(&validatorController)
	if query.RecordNotFound() {
		if queryController.RecordNotFound() {
			query = d.db.Create(&model.ValidatorInfo{ValidatorController: controller, ValidatorStash: stash})
		} else {
			query = d.db.Model(&queryController).Update(&model.ValidatorInfo{ValidatorStash: stash})
		}
	} else {
		if queryController.RecordNotFound() {
			query = d.db.Model(&validatorStash).Update(&model.ValidatorInfo{ValidatorController: controller})
		} else {
			d.db.Model(&validatorStash).Delete(&validatorStash)
			d.db.Model(&validatorController).Update(&model.ValidatorInfo{ValidatorStash: stash})
		}
	}
	return query.Error
}

func (d *Dao) SetRewardAccount(c context.Context, stash, rewardAccountType string) error {
	if stash == "" {
		return nil
	}
	var sqlQuery string
	if rewardAccountType == "Stash" {
		sqlQuery = fmt.Sprintf("Update validator_infos Set validator_reward = validator_stash Where validator_controller='%s'", stash)
	} else {
		sqlQuery = fmt.Sprintf("Update validator_infos Set validator_reward = validator_controller Where validator_controller='%s'", stash)
	}
	query := d.db.Exec(sqlQuery)
	return query.Error
}
