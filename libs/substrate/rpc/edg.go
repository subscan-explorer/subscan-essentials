package rpc

import (
	"github.com/itering/subscan/libs/substrate/storage"
)

type edg struct {
	query
}

func (r *edg) StakingStakers(stash string, currentEra int) *storage.Exposures {
	exposure, err := ReadStorage(nil, "Staking", "Stakers", "", stash)
	if err != nil {
		return nil
	}
	return exposure.ToExposures()
}
