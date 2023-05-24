package http

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	scale "github.com/itering/scale.go/types"
	"github.com/itering/scale.go/types/scaleBytes"
	"github.com/itering/subscan/plugins/router"
	"github.com/itering/subscan/plugins/staking/dao"
	"github.com/itering/subscan/plugins/staking/service"
	"github.com/itering/subscan/util"
	"github.com/itering/subscan/util/address"
	"github.com/itering/subscan/util/validator"
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
	w := c.Writer
	r := c.Request
	p := new(AddressReq)
	if err := validator.Validate(r.Body, p); err != nil {
		toJson(w, 10001, nil, err)
		return
	}
	depthConstant := svc.GetRuntimeConstant("Staking", "HistoryDepth")

	if depthConstant == nil {
		slog.Error("get runtime constant failed", "module", "Staking", "name", "HistoryDepth")
		toJson(w, 10001, nil, errors.New("get runtime constant failed"))
		return
	}

	m := scale.ScaleDecoder{}
	m.Init(scaleBytes.ScaleBytes{Data: util.HexToBytes(depthConstant.Value)}, nil)
	depth := m.ProcessAndUpdateData("U32").(uint32)

	activeEra := dao.GetLatestEra(svc.Storage())
	minEra := activeEra - depth

	list, _ := svc.GetPayoutListJson(p.Page, p.Row, p.Address, minEra)

	slog.Debug("RewardsSlashes", "page", p.Page, "row", p.Row, "address", p.Address, "found", len(list))

	toJson(w, 0, map[string]interface{}{
		"list": list, "count": len(list),
	}, nil)
}

// right now the staking dashboard only uses the era and number of points
type EraStat struct {
	Era         uint32 `json:"era"`
	RewardPoint uint32 `json:"reward_point"`
}

func eraStat(c *gin.Context) {
	w := c.Writer
	r := c.Request
	p := new(AddressReq)
	if err := validator.Validate(r.Body, p); err != nil {
		toJson(w, 10001, nil, err)
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

	toJson(w, 0, map[string]interface{}{
		"list": eraStats, "count": len(eraStats),
	}, nil)
}

func poolRewards(c *gin.Context) {
	w := c.Writer
	r := c.Request
	p := new(AddressReq)
	if err := validator.Validate(r.Body, p); err != nil {
		toJson(w, 10001, nil, err)
		return
	}

	addressSS58 := address.SS58Address(p.Address)
	list, _ := dao.GetPoolPayoutList(svc.Storage(), p.Page, p.Row, addressSS58)

	slog.Debug("PoolRewards", "page", p.Page, "row", p.Row, "address", p.Address, "found", len(list))

	toJson(w, 0, map[string]interface{}{
		"list": list, "count": len(list),
	}, nil)
}

type J struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	TTL     int         `json:"ttl"`
	Data    interface{} `json:"data,omitempty"`
}

func (j J) Render(w http.ResponseWriter) error {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = []string{"application/json; charset=utf-8"}
		header["Access-Control-Allow-Origin"] = []string{"*"}
	}
	return nil
}

func (j J) WriteContentType(w http.ResponseWriter) {
	var (
		jsonBytes []byte
		err       error
	)
	_ = j.Render(w)
	if jsonBytes, err = json.Marshal(j); err != nil {
		_ = errors.WithStack(err)
		return
	}
	if _, err = w.Write(jsonBytes); err != nil {
		_ = errors.WithStack(err)
	}
}

func toJson(w http.ResponseWriter, code int, data interface{}, err error) {
	j := J{
		Message: "success",
		TTL:     1,
		Data:    data,
	}
	if err != nil {
		j.Message = err.Error()
	}
	if code != 0 {
		j.Code = code
	}
	j.WriteContentType(w)
	_ = j.Render(w)
}
