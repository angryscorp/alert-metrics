package gzipper

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUnzipMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should return 400 for invalid gzip data", func(t *testing.T) {
		engine := gin.New()
		engine.Use(UnzipMiddleware())
		engine.POST("/test", func(c *gin.Context) {
			c.String(http.StatusOK, "should not reach here")
		})

		req := httptest.NewRequest(http.MethodPost, "/test", io.NopCloser(strings.NewReader("invalid gzip")))
		req.Header.Set("Content-Encoding", "gzip")

		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
