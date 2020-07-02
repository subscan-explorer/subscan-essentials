package http

import (
	"github.com/itering/subscan/internal/service"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/go-kratos/kratos/pkg/conf/paladin"
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	"github.com/stretchr/testify/assert"
)

func init() {
	_ = os.Setenv("TEST_MOD", "true")
	if client, err := paladin.NewFile("../../../configs"); err != nil {
		panic(err)
	} else {
		paladin.DefaultClient = client
	}
	svc = service.New()
	ss = svc.NewScan()
}

func testRequest(w *httptest.ResponseRecorder, req *http.Request) {
	e := bm.DefaultServer(nil)
	req.RemoteAddr = "127.0.0.1:8080"
	initRouter(e)
	e.ServeHTTP(w, req)
}

var testCases = []struct {
	url     string
	payload io.Reader
}{
	{"/api/scan/metadata", nil},
	{"/api/scan/blocks", strings.NewReader(`{"row": 10, "page": 0}`)},
	{"/api/scan/block", strings.NewReader(`{"block_hash": "0xbadc6963e1add4d7a588e350d837579491d08bb270f02c56b3dd5f17018dee0c"}`)},
	{"/api/scan/extrinsics", strings.NewReader(`{"row": 10, "page": 0}`)},
	{"/api/scan/extrinsic", strings.NewReader(`{"hash": "0xbadc6963e1add4d7a588e350d837579491d08bb270f02c56b3dd5f17018dee0c"}`)},
	{"/api/scan/events", strings.NewReader(`{"row": 10, "page": 0}`)},
	{"/api/scan/check_hash", strings.NewReader(`{"hash": "0xbadc6963e1add4d7a588e350d837579491d08bb270f02c56b3dd5f17018dee0c"}`)},
	{"/api/scan/runtime/metadata", strings.NewReader(`{"spec": 1}`)},
	{"/api/scan/runtime/list", nil},
}

func TestRouter(t *testing.T) {
	for _, test := range testCases {
		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", test.url, test.payload)
		assert.NotNil(t, req)

		req.Header.Set("Content-Type", "application/json")
		testRequest(w, req)

		assert.NoError(t, err)
		assert.Equal(t, 200, w.Code)
	}
}
