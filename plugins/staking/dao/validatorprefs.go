package dao

import (
	"encoding/json"
	"errors"

	"github.com/itering/subscan-plugin/storage"
	"github.com/itering/subscan/plugins/staking/model"
	"github.com/shopspring/decimal"
	"golang.org/x/exp/slog"
)

type CommissionHistoryRecord struct {
	StartBlock uint32          `json:"startBlock"`
	NewValue   decimal.Decimal `json:"newValue"`
}

type CommissionHistory []CommissionHistoryRecord

func NewValidatorPrefs(db storage.DB, addressSS58 string, commission decimal.Decimal, blockedNomination bool, blockNumber uint32) error {
	slog.Info("NewValidatorPrefs", "account", addressSS58, "commission", commission, "blockedNomination", blockedNomination)
	var maybe []model.ValidatorPrefs
	opt := storage.Option{PluginPrefix: "staking"}
	db.FindBy(&maybe, map[string]interface{}{"account": addressSS58}, &opt)
	if len(maybe) == 1 {
		prefs := maybe[0]
		prefs.Commission = commission
		prefs.BlockedNomination = blockedNomination
		var history CommissionHistory
		if err := json.Unmarshal([]byte(prefs.CommissionHistory), &history); err != nil {
			return err
		}
		history = append(history, CommissionHistoryRecord{
			StartBlock: blockNumber,
			NewValue:   commission,
		})
		ch, err := json.Marshal(history)
		if err != nil {
			return err
		}
		prefs.CommissionHistory = string(ch)
		if err := db.Update(&prefs, map[string]interface{}{"ID": prefs.ID}, map[string]interface{}{
			"commission":         prefs.Commission,
			"blocked_nomination": prefs.BlockedNomination,
			"commission_history": prefs.CommissionHistory,
		}); err != nil {
			return err
		}
	} else {
		var history CommissionHistory
		history = append(history, CommissionHistoryRecord{
			StartBlock: blockNumber,
			NewValue:   commission,
		})
		commissionHistory, err := json.Marshal(history)
		if err != nil {
			return err
		}

		if err := db.Create(&model.ValidatorPrefs{Account: addressSS58, Commission: commission, BlockedNomination: blockedNomination, CommissionHistory: string(commissionHistory)}); err != nil {
			return err
		}
	}

	return nil
}

var ErrNotFound = errors.New("model not found")

func GetValidatorPrefs(db storage.DB, validatorAddressSS58 string) (model.ValidatorPrefs, error) {
	var prefs model.ValidatorPrefs
	opt := storage.Option{PluginPrefix: "staking", PageSize: 1}
	if count, _ := db.FindBy(&prefs, map[string]interface{}{"account": validatorAddressSS58}, &opt); count == 0 {
		return prefs, ErrNotFound
	}
	return prefs, nil
}
