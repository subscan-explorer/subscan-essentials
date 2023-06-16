package observer

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/itering/subscan/configs"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	endpoints := []string{"", "ping", "health"}
	config := &configs.HealthCheck{
		Addr: ":80",
	}
	router := newHealthCheckServer(config)

	for _, endpoint := range endpoints {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/"+endpoint, nil)
		router.Handler.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "ok", w.Body.String())
	}
}
