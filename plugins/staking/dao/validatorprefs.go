package dao

import (
	"fmt"

	"github.com/itering/subscan/plugins/staking/model"
	"github.com/itering/subscan/plugins/storage"
	"github.com/itering/subscan/util/address"
	"github.com/shopspring/decimal"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
)

type CommissionHistoryRecord struct {
	StartBlock uint32          `json:"startBlock"`
	NewValue   decimal.Decimal `json:"newValue"`
}

type CommissionHistory []CommissionHistoryRecord

func NewValidatorPrefs(db storage.DB, addressSS58 address.SS58Address, commission decimal.Decimal, blockedNomination bool, blockNumber uint32) error {
	var info model.EraInfo
	res := db.Query(&model.EraInfo{}).Select("era, start_block").Where("start_block <= ?", blockNumber).Order("start_block DESC").Limit(1).Find(&info)
	if res.Error != nil {
		slog.Error("NewValidatorPrefs", "account", addressSS58, "commission", commission, "blockedNomination", blockedNomination, "blockNumber", blockNumber, "error", res.Error)
		return res.Error
	}
	slog.Warn("NewValidatorPrefs", "account", addressSS58, "commission", commission, "blockedNomination", blockedNomination, "era", info.Era, "blockNumber", blockNumber)
	var maybe []model.ValidatorPrefs
	opt := storage.Option{PluginPrefix: "staking"}
	db.FindBy(&maybe, map[string]interface{}{"account": addressSS58}, &opt)
	if len(maybe) > 0 {
		slog.Warn("NewValidatorPrefs", "account", addressSS58, "commission", commission, "blockedNomination", blockedNomination, "era", info.Era, "maybe", maybe[0])
		existing := maybe[0]
		existing.Commission = commission
		existing.BlockedNomination = blockedNomination
		return db.Query(&model.ValidatorPrefs{}).Save(&existing).Error
	}
	if err := db.Create(&model.ValidatorPrefs{Account: addressSS58, Commission: commission, BlockedNomination: blockedNomination, Era: info.Era}); err != nil {
		return err
	}

	return nil
}

func GetValidatorPrefs(db storage.DB, validatorAddressSS58 string, era uint32) (model.ValidatorPrefs, error) {
	var prefs model.ValidatorPrefs
	res := db.Query(&prefs).Select("*").Where("account = ? AND era <= ?", validatorAddressSS58, era).Order("era DESC").Limit(1).Find(&prefs)
	if res.Error == gorm.ErrRecordNotFound {
		return prefs, fmt.Errorf("validator prefs model not found for %s %d", validatorAddressSS58, era)
	}
	return prefs, nil
}
