package http

import (
	"github.com/gin-gonic/gin"
	scale "github.com/itering/scale.go/types"
	"github.com/itering/scale.go/types/scaleBytes"
	"github.com/itering/subscan/plugins/router"
	"github.com/itering/subscan/plugins/staking/dao"
	"github.com/itering/subscan/plugins/staking/service"
	"github.com/itering/subscan/util"
	"github.com/itering/subscan/util/address"
	"github.com/pkg/errors"
	"golang.org/x/exp/slog"
)

var svc *service.Service

func Router(s *service.Service) []router.Http {
	svc = s
	return []router.Http{{Router: "eraStat", Handle: eraStat}, {Router: "rewardsSlashes", Handle: rewardsSlashes}, {Router: "poolRewards", Handle: poolRewards}}
}

type AddressReq struct {
	Row     int    `json:"row" binding:"min=1,max=5000"`
	Page    int    `json:"page" binding:"min=0"`
	Address string `json:"address" binding:"required"`
}

func rewardsSlashes(c *gin.Context) {
	p := new(AddressReq)
	if err := c.BindJSON(p); err != nil {
		util.ToJson(c, nil, util.ParamsError)
	}
	depthConstant := svc.GetRuntimeConstant("Staking", "HistoryDepth")

	if depthConstant == nil {
		slog.Error("get runtime constant failed", "module", "Staking", "name", "HistoryDepth")
		util.ToJson(c, nil, errors.New("get runtime constant failed"))
		return
	}

	m := scale.ScaleDecoder{}
	m.Init(scaleBytes.ScaleBytes{Data: util.HexToBytes(depthConstant.Value)}, nil)
	depth := m.ProcessAndUpdateData("U32").(uint32)

	activeEra := dao.GetLatestEra(svc.Storage())
	minEra := activeEra - depth

	list, _ := svc.GetPayoutListJson(p.Page, p.Row, p.Address, minEra)

	slog.Debug("RewardsSlashes", "page", p.Page, "row", p.Row, "address", p.Address, "found", len(list))

	util.ToJson(c, map[string]interface{}{
		"list": list, "count": len(list),
	}, nil)
}

// right now the staking dashboard only uses the era and number of points
type EraStat struct {
	Era         uint32 `json:"era"`
	RewardPoint uint32 `json:"reward_point"`
}

func eraStat(c *gin.Context) {
	p := new(AddressReq)
	if err := c.BindJSON(p); err != nil {
		util.ToJson(c, nil, err)
		return
	}

	addressSS58 := address.SS58Address(p.Address)

	eraAndPointsList, _ := dao.GetEraPointsList(svc.Storage(), p.Page, p.Row)

	eraStats := make([]EraStat, 0, len(eraAndPointsList))
	for _, eraAndPoints := range eraAndPointsList {
		eraStats = append(eraStats, EraStat{
			Era:         eraAndPoints.Era,
			RewardPoint: eraAndPoints.Points[addressSS58],
		})
	}

	slog.Debug("EraStat", "page", p.Page, "row", p.Row, "address", p.Address, "found", len(eraStats))

	util.ToJson(c, map[string]interface{}{
		"list": eraStats, "count": len(eraStats),
	}, nil)
}

func poolRewards(c *gin.Context) {
	p := new(AddressReq)
	if err := c.BindJSON(p); err != nil {
		util.ToJson(c, nil, err)
		return
	}

	addressSS58 := address.SS58Address(p.Address)
	list, _ := dao.GetPoolPayoutList(svc.Storage(), p.Page, p.Row, addressSS58)

	slog.Debug("PoolRewards", "page", p.Page, "row", p.Row, "address", p.Address, "found", len(list))

	util.ToJson(c, map[string]interface{}{
		"list": list, "count": len(list),
	}, nil)
}
