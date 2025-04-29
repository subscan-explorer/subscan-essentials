package http

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
)

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

type EtherRes struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

func (j EtherRes) Render(w http.ResponseWriter) error {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = []string{"application/json; charset=utf-8"}
	}
	return nil
}

func (j EtherRes) WriteContentType(w http.ResponseWriter) {
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

func etherscanRes(w http.ResponseWriter, code int, data interface{}, messageErr error) {
	j := EtherRes{
		Message: "ok",
		Result:  data,
	}
	if messageErr != nil {
		j.Message = messageErr.Error()
	}
	j.Status = code
	if data == nil {
		j.Status = 0
		j.Message = "No records found"
	}
	j.WriteContentType(w)
	_ = j.Render(w)
}
