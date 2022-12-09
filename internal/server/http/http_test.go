package http

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/itering/subscan/configs"
	"github.com/itering/subscan/internal/service"
	"github.com/itering/subscan/util"

	"github.com/stretchr/testify/assert"
)

func init() {
	util.ConfDir = "../../../configs"
	configs.Init()
	svc = service.New()
}

func testRequest(w *httptest.ResponseRecorder, req *http.Request) {
	e := gin.New()
	gin.SetMode(gin.ReleaseMode)
	req.RemoteAddr = "127.0.0.1:8080"
	initRouter(e)
	e.ServeHTTP(w, req)
}

var testCases = []struct {
	url     string
	payload io.Reader
	method  string
}{
	{"/api/scan/metadata", nil, "POST"},
	{"/api/scan/blocks", strings.NewReader(`{"row": 10, "page": 0}`), "POST"},
	{"/api/scan/block", strings.NewReader(`{"block_hash": "0xbadc6963e1add4d7a588e350d837579491d08bb270f02c56b3dd5f17018dee0c"}`), "POST"},
	{"/api/scan/extrinsics", strings.NewReader(`{"row": 10, "page": 0}`), "POST"},
	{"/api/scan/extrinsic", strings.NewReader(`{"hash": "0xbadc6963e1add4d7a588e350d837579491d08bb270f02c56b3dd5f17018dee0c"}`), "POST"},
	{"/api/scan/events", strings.NewReader(`{"row": 10, "page": 0}`), "POST"},
	{"/api/scan/check_hash", strings.NewReader(`{"hash": "0xbadc6963e1add4d7a588e350d837579491d08bb270f02c56b3dd5f17018dee0c"}`), "POST"},
	{"/api/scan/runtime/metadata", strings.NewReader(`{"spec": 1}`), "POST"},
	{"/api/scan/runtime/list", nil, "POST"},
	{"/api/now", nil, "POST"},
	{"/api/system/status", nil, "GET"},
	{"/ping", nil, "GET"},
}

func TestRouter(t *testing.T) {
	for _, test := range testCases {
		w := httptest.NewRecorder()
		req, err := http.NewRequest(test.method, test.url, test.payload)
		assert.NotNil(t, req)

		req.Header.Set("Content-Type", "application/json")
		testRequest(w, req)

		assert.NoError(t, err)
		assert.Equal(t, 200, w.Code)
	}
}
