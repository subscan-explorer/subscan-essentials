package tasks

import (
	"context"
	"fmt"

	"github.com/itering/subscan/internal/service"
	"github.com/itering/subscan/util"
)

func RefreshMetadata(srv *service.Service) {
	c := context.TODO()
	u := map[string]interface{}{
		"count_account":           srv.GetActiveAccountCount(c),
		"count_signed_extrinsic":  srv.GetTransactionCount(c),
		"current_validator_count": util.IntToString(srv.CurrentValidatorsCount(c)), // 当前 Validator 的数目
	}
	if validatorCount, err := srv.GetValidatorCount(c); err != nil || validatorCount == 0 {
		u["validator_count"] = u["current_validator_count"]
	} else {
		u["validator_count"] = util.IntToString(validatorCount) // 节点允许 validator 的数目
	}
	u["waiting_validator"] = srv.GetWaitingValidatorCount()
	u["count_transfer"] = srv.GetTransferCount()

	_ = srv.UpdateChainMetadata(u)
	fmt.Println("finish refresh metadata")
}
