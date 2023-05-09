package dao

import (
	"github.com/itering/subscan-plugin/storage"
	scanModel "github.com/itering/subscan/model"
	"github.com/itering/subscan/plugins/staking/model"
	"github.com/itering/subscan/util/address"
	"github.com/shopspring/decimal"
	"golang.org/x/exp/slog"
)

func NewClaimedPayout(db storage.DB, addressHex string, validatorAccountSs58 string, amount decimal.Decimal, era uint32, event *scanModel.ChainEvent, block *scanModel.ChainBlock, extrinsicIndex string) error {
	accountId := address.SS58Address(addressHex)
	slog.Info("NewClaimedPayout", "account", accountId, "validator", validatorAccountSs58, "amount", amount)
	opt := storage.Option{PluginPrefix: "staking"}
	var unclaimedPayout []model.Payout
	db.FindBy(&unclaimedPayout, map[string]interface{}{"account": accountId, "era": era}, &opt)
	if len(unclaimedPayout) == 1 {
		for _, payout := range unclaimedPayout {
			if payout.Amount.Equal(amount) {
				if err := db.Update(&payout, map[string]interface{}{"claimed": true}, &opt); err != nil {
					return err
				}
			} else {
				slog.Error("Found unexpected amount of unclaimed payouts", "unclaimedPayout", unclaimedPayout, "amount", amount)
			}
		}
	} else {
		slog.Error("Found unexpected number of unclaimed payouts", "unclaimedPayout", unclaimedPayout, "amount", amount, "era", era, "account", accountId)
	}

	if err := db.Create(&model.Payout{
		Account:        accountId,
		Amount:         amount,
		BlockTimestamp: uint64(block.BlockTimestamp),
		Stash:          accountId,
		ValidatorStash: validatorAccountSs58,
		ModuleId:       "staking",
		EventId:        "Rewarded",
		ExtrinsicIndex: extrinsicIndex,
		Era:            era,
		EventIndex:     event.EventIndex,
	}); err != nil {
		return err
	}
	return nil
}

func NewUnclaimedPayout(db storage.DB, addressHex string, validatorAccountSs58 string, amount decimal.Decimal, era uint32) error {
	accountId := address.SS58Address(addressHex)
	slog.Info("NewUnclaimedPayout", "account", accountId, "validator", validatorAccountSs58, "amount", amount)
	_ = db.Create(&model.Payout{Account: accountId, Amount: amount, Stash: accountId, ValidatorStash: validatorAccountSs58, Era: era})
	return nil
}

func AfterClaimedPayoutCreate(db storage.DB, accountId string) error {
	// accountDataRaw, err := rpc.ReadStorage(nil, "staking", "nominations", "", accountId)
	// if err != nil {
	// 	return err
	// }
	// slog.Debug("accountDataRaw: ", accountDataRaw)
	return nil
}

func GetPayoutList(db storage.DB, page, row int, address string) ([]model.Payout, int) {
	var claimedPayouts []model.Payout
	opt := storage.Option{PluginPrefix: "staking", Page: page, PageSize: row}
	db.FindBy(&claimedPayouts, map[string]interface{}{"account": address}, &opt)
	return claimedPayouts, len(claimedPayouts)
}
