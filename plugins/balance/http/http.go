package http

import (
	"encoding/json"
	"github.com/itering/subscan-plugin/router"
	"github.com/itering/subscan/plugins/balance/service"
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
		{"accounts", accounts},
	}
}

func accounts(w http.ResponseWriter, r *http.Request) error {
	p := new(struct {
		Row  int `json:"row" validate:"min=1,max=100"`
		Page int `json:"page" validate:"min=0"`
	})
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
