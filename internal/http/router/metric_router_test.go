package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestMetricRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := zerolog.Nop()

	t.Run("New router creation", func(t *testing.T) {
		engine := gin.New()
		metricRouter := New(engine, &logger)

		assert.NotNil(t, metricRouter)
		assert.Equal(t, engine, metricRouter.engine)
	})

	t.Run("Router has required methods", func(t *testing.T) {
		engine := gin.New()
		metricRouter := New(engine, &logger)

		assert.NotNil(t, metricRouter.RegisterPingHandler)
		assert.NotNil(t, metricRouter.RegisterMetricsHandler)
		assert.NotNil(t, metricRouter.RegisterMetricsJSONHandler)
		assert.NotNil(t, metricRouter.Run)
	})

	t.Run("Router registers routes correctly", func(t *testing.T) {
		engine := gin.New()
		metricRouter := New(engine, &logger)

		metricRouter.engine.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"test": "ok"})
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "test")
	})
}
