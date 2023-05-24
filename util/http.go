package util

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type JsonResult struct {
	Code        int         `json:"code"`
	Message     string      `json:"message"`
	GeneratedAt int64       `json:"generated_at"`
	Data        interface{} `json:"data,omitempty"`
}

type J struct {
	Code        int         `json:"code"`
	Message     string      `json:"message"`
	GeneratedAt int64       `json:"generated_at"`
	Data        interface{} `json:"data,omitempty"`
}

var jsonContentType = []string{"application/json; charset=utf-8"}

// WriteContentType (JSON) writes JSON ContentType.
func (r JsonResult) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}

// WriteJSON marshals the given interface object and writes it with custom ContentType.
func WriteJSON(w http.ResponseWriter, obj interface{}) error {
	writeContentType(w, jsonContentType)
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	_, err = w.Write(jsonBytes)
	return err
}

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}

// Render (JSON) writes data with custom ContentType.
func (r JsonResult) Render(w http.ResponseWriter) (err error) {
	if err = WriteJSON(w, r); err != nil {
		panic(err)
	}
	return
}

func ToJson(c *gin.Context, data interface{}, err error) {
	if ctxErr := c.Request.Context().Err(); ctxErr != nil && ctxErr == context.DeadlineExceeded {
		c.AbortWithStatus(500)
		c.Render(0, JsonResult{Code: 500, GeneratedAt: time.Now().Unix(), Message: ctxErr.Error()})
		return
	}
	j := JsonResult{
		Data:        data,
		GeneratedAt: time.Now().Unix(),
		Message:     "Success",
	}
	if err != nil {
		if ec, ok := errors.Cause(err).(ErrorCode); ok {
			j.Code = ec.Code()
			j.Message = ec.Message()
		} else {
			j.Code = 400
			j.Message = err.Error()
		}
	}
	c.Render(0, j)
}
