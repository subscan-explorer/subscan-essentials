package http

import (
	"encoding/json"
	"net/http"

	scale "github.com/itering/scale.go/types"
	"github.com/itering/scale.go/types/scaleBytes"
	"github.com/itering/subscan-plugin/router"
	"github.com/itering/subscan/plugins/staking/dao"
	"github.com/itering/subscan/plugins/staking/model"
	"github.com/itering/subscan/plugins/staking/service"
	"github.com/itering/subscan/util"
	"github.com/itering/subscan/util/validator"
	"github.com/pkg/errors"
	"golang.org/x/exp/slog"
)

var svc *service.Service

func Router(s *service.Service) []router.Http {
	svc = s
	return []router.Http{{Router: "rewardsSlashes", Handle: rewardsSlashes}}
}

func rewardsSlashes(w http.ResponseWriter, r *http.Request) error {
	p := new(struct {
		Row     int    `json:"row" validate:"min=1,max=5000"`
		Page    int    `json:"page" validate:"min=0"`
		Address string `json:"address" validate:"required"`
	})
	if err := validator.Validate(r.Body, p); err != nil {
		toJson(w, 10001, nil, err)
		return nil
	}

	list, _ := svc.GetPayoutListJson(p.Page, p.Row, p.Address)

	depthConstant := svc.Dao().GetRuntimeConstantLatest("Staking", "HistoryDepth")

	if depthConstant == nil {
		slog.Error("get runtime constant failed", "module", "Staking", "name", "HistoryDepth")
		toJson(w, 10001, nil, errors.New("get runtime constant failed"))
		return nil
	}

	m := scale.ScaleDecoder{}
	m.Init(scaleBytes.ScaleBytes{Data: util.HexToBytes(depthConstant.Value)}, nil)
	depth := m.ProcessAndUpdateData("U32").(uint32)

	activeEra := dao.GetLatestEra(svc.Storage())

	filteredList := make([]model.Payout, 0)

	for _, item := range list {
		if uint32(item.BlockTimestamp) == 0 && item.Era < activeEra-depth {
			continue
		}
		filteredList = append(filteredList, item)
	}

	slog.Debug("RewardsSlashes", "page", p.Page, "row", p.Row, "address", p.Address, "found", len(filteredList))

	toJson(w, 0, map[string]interface{}{
		"list": filteredList, "count": len(filteredList),
	}, nil)
	return nil
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
