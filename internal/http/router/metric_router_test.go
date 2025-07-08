package router

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/angryscorp/alert-metrics/internal/domain"
	"github.com/angryscorp/alert-metrics/internal/infrastructure/metricstorage"
)

func TestMetricRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	store := metricstorage.NewMemoryMetricStorage()
	val := 123.456
	_ = store.UpdateMetric(context.TODO(), domain.Metric{
		ID:    "test123",
		MType: domain.MetricTypeGauge,
		Value: &val,
	})
	metricRouter := New(gin.New(), store)

	tests := []struct {
		name     string
		method   string
		path     string
		response int
	}{
		{
			name:     "Healthcheck returns StatusOK in case of no errors",
			method:   http.MethodGet,
			path:     "/ping",
			response: http.StatusOK,
		},
		{
			name:     "Update returns StatusOK in case of no errors",
			method:   http.MethodPost,
			path:     "/update/gauge/test/123",
			response: http.StatusOK,
		},
		{
			name:     "Update returns StatusBadRequest in case of wrong metric type",
			method:   http.MethodPost,
			path:     "/update/gauge2/test/123",
			response: http.StatusBadRequest,
		},
		{
			name:     "Update returns StatusNotFound in case of wrong method",
			method:   http.MethodGet,
			path:     "/update/gauge/test/123",
			response: http.StatusNotFound,
		},
		{
			name:     "Update returns StatusNotFound in case of wrong path",
			method:   http.MethodPost,
			path:     "/record/gauge/test/123",
			response: http.StatusNotFound,
		},
		{
			name:     "Update returns StatusBadRequest in case of wrong value",
			method:   http.MethodPost,
			path:     "/update/counter/test/123.456",
			response: http.StatusBadRequest,
		},
		{
			name:     "Value returns StatusNotFound in case of unknown metric name",
			method:   http.MethodGet,
			path:     "/value/counter/test",
			response: http.StatusNotFound,
		},
		{
			name:     "Value returns StatusOK in case of metric name is exist",
			method:   http.MethodGet,
			path:     "/value/gauge/test123",
			response: http.StatusOK,
		},
		{
			name:     "Value returns StatusNotFound in case of wrong method",
			method:   http.MethodPost,
			path:     "/value/gauge/test123",
			response: http.StatusNotFound,
		},
		{
			name:     "Home page returns StatusOK in case of no errors",
			method:   http.MethodGet,
			path:     "/",
			response: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			request := httptest.NewRequest(test.method, test.path, nil)
			w := httptest.NewRecorder()

			// Act
			metricRouter.engine.ServeHTTP(w, request)
			res := w.Result()
			_ = res.Body.Close()

			// Assert
			require.NotNil(t, res)
			require.Equal(t, test.response, res.StatusCode)
		})
	}
}
