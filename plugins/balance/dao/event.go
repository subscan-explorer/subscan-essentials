package dao

import (
	"context"
	"github.com/itering/subscan-plugin/storage"
)

func EmitEvent(ctx context.Context, d storage.Dao, event *storage.Event) error {
	switch event.EventId {
	// [accountId, balance]
	case "Endowed", "Reserved", "Unreserved", "Deposit", "Minted", "Issued", "Locked", "Unlocked", "Withdraw":

	}
	return nil
}
