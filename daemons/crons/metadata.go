package crons

import (
	"context"
	"fmt"
	"github.com/itering/subscan/internal/service"
)

func RefreshMetadata(srv *service.Service) {
	c := context.TODO()
	u := map[string]interface{}{
		"count_signed_extrinsic": srv.GetTransactionCount(c),
	}
	_ = srv.UpdateChainMetadata(u)
	fmt.Println("finish refresh metadata")
}
