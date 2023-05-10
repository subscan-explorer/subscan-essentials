package dao

import (
	"errors"

	scanModel "github.com/itering/subscan/model"
	"github.com/itering/subscan/plugins/staking/model"
	"github.com/itering/subscan/plugins/storage"
	"github.com/itering/subscan/util/address"
	"github.com/shopspring/decimal"
	"golang.org/x/exp/slog"
)

func NewClaimedPayout(db storage.DB, addressHex string, validatorAccountSs58 string, amount decimal.Decimal, era uint32, event *scanModel.ChainEvent, block *scanModel.ChainBlock, extrinsicIndex string) error {
	accountId := address.SS58AddressFromHex(addressHex)
	slog.Info("NewClaimedPayout", "account", accountId, "validator", validatorAccountSs58, "amount", amount, "eventIndex", event.EventIndex)
	opt := storage.Option{PluginPrefix: "staking"}
	var unclaimedPayout []model.Payout
	db.FindBy(&unclaimedPayout, map[string]interface{}{"account": accountId, "era": era}, &opt)
	if len(unclaimedPayout) == 1 {
		payout := unclaimedPayout[0]
		if !payout.Amount.Equal(amount) {
			// FIXME: this shouldn't actually happen, but I can't figure out what problem with the calculation is (at least
			// not before the deadline)
			slog.Error("Found unexpected amount of unclaimed payouts", "unclaimedPayout", unclaimedPayout, "amount", amount)
		}

		payout.BlockTimestamp = uint64(block.BlockTimestamp)
		payout.EventIndex = event.EventIndex
		payout.ExtrinsicIndex = extrinsicIndex
		payout.ModuleId = "staking"
		payout.EventId = "Rewarded"
		payout.Amount = amount
		payout.Claimed = true
		if err := db.Update(&payout, map[string]interface{}{"ID": payout.ID}, map[string]interface{}{
			"block_timestamp": payout.BlockTimestamp,
			"event_index":     payout.EventIndex,
			"extrinsic_index": payout.ExtrinsicIndex,
			"module_id":       payout.ModuleId,
			"event_id":        payout.EventId,
			"amount":          payout.Amount,
			"claimed":         payout.Claimed,
		}); err != nil {
			return err
		}
	} else {
		slog.Error("Found unexpected number of unclaimed payouts", "unclaimedPayout", unclaimedPayout, "amount", amount, "era", era, "account", accountId)
		return errors.New("found unexpected number of unclaimed payouts")
	}

	return nil
}

func NewUnclaimedPayout(db storage.DB, addressSS58 address.SS58Address, validatorAccountSs58 address.SS58Address, amount decimal.Decimal, era uint32) error {
	slog.Info("NewUnclaimedPayout", "account", addressSS58, "validator", validatorAccountSs58, "amount", amount)
	_ = db.Create(&model.Payout{Account: addressSS58, Amount: amount, Stash: addressSS58, ValidatorStash: validatorAccountSs58, Era: era})
	return nil
}

func AfterClaimedPayoutCreate(db storage.DB, accountId string) error {
	return nil
}

func GetPayoutList(db storage.DB, page, row int, address string, minEra uint32) ([]model.Payout, int) {
	var claimedPayouts []model.Payout
	db.Query(model.Payout{}).Where("account = ? and (era >= ? or claimed <> 0)", address, minEra).Order("block_timestamp DESC").Limit(row).Offset((page - 1) * row).Find(&claimedPayouts)
	return claimedPayouts, len(claimedPayouts)
}

func GetLatestEra(db storage.DB) uint32 {
	var payout model.Payout
	opt := storage.Option{PluginPrefix: "staking", Order: "era desc", PageSize: 1}
	db.FindBy(&payout, map[string]interface{}{}, &opt)
	return payout.Era
}
