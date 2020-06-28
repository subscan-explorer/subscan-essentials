package middlewares

import (
	"fmt"
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
	time2 "github.com/go-kratos/kratos/pkg/time"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// CORS middleware
// Set header with Access-Control-*
func Test_CORS(t *testing.T) {
	var server bm.ServerConfig
	server.Network = "0.0.0.0:4399"
	server.Timeout = time2.Duration(time.Second)
	engine := bm.DefaultServer(&server)
	engine.HandleMethodNotAllowed = false
	engine.Use(CORS())
	engine.GET("/", func(c *bm.Context) {})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/", nil)
	req.RemoteAddr = "127.0.0.1:8080"
	assert.NotNil(t, req)

	req.Header.Set("Content-Type", "application/json")

	engine.ServeHTTP(w, req)

	fmt.Println(w.Header())
	assert.Equal(t, w.Header(), http.Header(
		http.Header{"Access-Control-Allow-Credentials": []string{"true"},
			"Access-Control-Allow-Headers":  []string{"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With"},
			"Access-Control-Allow-Methods":  []string{"POST, GET, OPTIONS, PUT, DELETE"},
			"Access-Control-Allow-Origin":   []string{"*"},
			"Access-Control-Expose-Headers": []string{"Content-Length"},
			"Access-Control-Max-Age":        []string{"86400"},
			"Content-Type":                  []string{"text/plain"},
			"Kratos-Trace-Id":               []string{""}}),
	)
}
