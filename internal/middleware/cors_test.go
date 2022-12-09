package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// CORS middleware
// Set header with Access-Control-*
func Test_CORS(t *testing.T) {
	engine := gin.New()
	engine.HandleMethodNotAllowed = false
	engine.Use(CORS())
	engine.GET("/", func(c *gin.Context) {})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/", nil)
	req.RemoteAddr = "127.0.0.1:8080"
	assert.NotNil(t, req)
	req.Header.Set("Content-Type", "application/json")
	engine.ServeHTTP(w, req)
	assert.Equal(t, w.Header()["Access-Control-Allow-Origin"], []string{"*"})
	assert.Equal(t, w.Header()["Access-Control-Allow-Methods"], []string{"POST, GET, OPTIONS, PUT, DELETE"})
}
