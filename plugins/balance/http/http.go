package http

import (
	"encoding/json"
	"github.com/itering/subscan-plugin/router"
	_ "github.com/itering/subscan/plugins/balance/model"
	"github.com/itering/subscan/plugins/balance/service"
	"github.com/itering/subscan/util/address"
	"github.com/itering/subscan/util/validator"
	"github.com/pkg/errors"
	"net/http"
)

var (
	svc *service.Service
)

func Router(s *service.Service) []router.Http {
	svc = s
	return []router.Http{
		{"accounts", accountsHandle, http.MethodPost},
		{"account", accountHandle, http.MethodPost},
		{"transfer", transferHandle, http.MethodPost},
	}
}

type accountsParams struct {
	Row  int `json:"row" validate:"min=1,max=100"`
	Page int `json:"page" validate:"min=0"`
}

// @Summary Get accounts list
// @Tags accounts
// @Accept json
// @Produce json
// @Param params body accountsParams true "params"
// @Success 200 {object} J{data=object{list=[]model.Account,count=int}}
// @Router /api/plugin/balance/accounts [post]
func accountsHandle(w http.ResponseWriter, r *http.Request) error {
	p := new(accountsParams)
	if err := validator.Validate(r.Body, p); err != nil {
		toJson(w, 10001, nil, err)
		return nil
	}

	list, count := svc.GetAccountListJson(p.Page, p.Row)
	toJson(w, 0, map[string]interface{}{
		"list": list, "count": count,
	}, nil)
	return nil
}

type accountParams struct {
	Address string `json:"address" validate:"required,addr"`
}

// @Summary Get account details
// @Tags accounts
// @Accept json
// @Produce json
// @Param params body accountParams true "params"
// @Success 200 {object} J{data=model.Account}
// @Router /api/plugin/balance/account [post]
func accountHandle(w http.ResponseWriter, r *http.Request) error {
	p := new(accountParams)
	if err := validator.Validate(r.Body, p); err != nil {
		toJson(w, 10001, nil, err)
		return nil
	}
	account := svc.GetAccountJson(r.Context(), address.Decode(p.Address))
	toJson(w, 0, account, nil)
	return nil
}

type transferParams struct {
	Address  string `json:"address" validate:"omitempty,addr"`
	BlockNum uint   `json:"block_num" validate:"omitempty,min=0"`
	Row      int    `json:"row" validate:"min=1,max=100"`
	Page     int    `json:"page" validate:"min=0"`
}

// @Summary Get transfer list
// @Tags transfers
// @Accept json
// @Produce json
// @Param params body transferParams true "params"
// @Success 200 {object} J{data=object{list=[]model.Transfer,count=int}}
// @Router /api/plugin/balance/transfer [post]
func transferHandle(w http.ResponseWriter, r *http.Request) error {
	p := new(transferParams)
	if err := validator.Validate(r.Body, p); err != nil {
		toJson(w, 10001, nil, err)
		return nil
	}
	list, count := svc.GetTransferJson(r.Context(), address.Decode(p.Address), p.BlockNum, p.Page, p.Row)
	toJson(w, 0, map[string]interface{}{
		"list": list, "count": count,
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
