package handler

import (
	"errors"
	"github.com/angryscorp/alert-metrics/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMetricsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name      string
		method    string
		path      string
		response  int
		setupMock func(*MockMetricStorage)
	}{
		{
			name:     "GetMetric returns StatusOK for existing counter",
			method:   http.MethodGet,
			path:     "/value/counter/test_counter",
			response: http.StatusOK,
			setupMock: func(m *MockMetricStorage) {
				m.On("GetMetric", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("domain.MetricType"), "test_counter").
					Return(domain.Metric{
						ID:    "test_counter",
						MType: "counter",
						Delta: func() *int64 { v := int64(42); return &v }(),
					}, true)
			},
		},
		{
			name:     "GetMetric returns StatusOK for existing gauge",
			method:   http.MethodGet,
			path:     "/value/gauge/test_gauge",
			response: http.StatusOK,
			setupMock: func(m *MockMetricStorage) {
				m.On("GetMetric", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("domain.MetricType"), "test_gauge").
					Return(domain.Metric{
						ID:    "test_gauge",
						MType: "gauge",
						Value: func() *float64 { v := 3.14; return &v }(),
					}, true)
			},
		},
		{
			name:     "GetMetric returns StatusNotFound for non-existing metric",
			method:   http.MethodGet,
			path:     "/value/counter/unknown",
			response: http.StatusNotFound,
			setupMock: func(m *MockMetricStorage) {
				m.On("GetMetric", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("domain.MetricType"), "unknown").
					Return(domain.Metric{}, false)
			},
		},
		{
			name:      "GetMetric returns StatusBadRequest for invalid metric type",
			method:    http.MethodGet,
			path:      "/value/invalid/test",
			response:  http.StatusBadRequest,
			setupMock: func(m *MockMetricStorage) {},
		},
		{
			name:     "UpdateMetrics returns StatusOK for valid counter",
			method:   http.MethodPost,
			path:     "/update/counter/test/123",
			response: http.StatusOK,
			setupMock: func(m *MockMetricStorage) {
				m.On("UpdateMetric", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("domain.Metric")).
					Return(nil)
			},
		},
		{
			name:     "UpdateMetrics returns StatusOK for valid gauge",
			method:   http.MethodPost,
			path:     "/update/gauge/test/123.456",
			response: http.StatusOK,
			setupMock: func(m *MockMetricStorage) {
				m.On("UpdateMetric", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("domain.Metric")).
					Return(nil)
			},
		},
		{
			name:      "UpdateMetrics returns StatusBadRequest for invalid metric type",
			method:    http.MethodPost,
			path:      "/update/invalid/test/123",
			response:  http.StatusBadRequest,
			setupMock: func(m *MockMetricStorage) {},
		},
		{
			name:      "UpdateMetrics returns StatusBadRequest for invalid counter value",
			method:    http.MethodPost,
			path:      "/update/counter/test/123.456",
			response:  http.StatusBadRequest,
			setupMock: func(m *MockMetricStorage) {},
		},
		{
			name:      "UpdateMetrics returns StatusBadRequest for invalid gauge value",
			method:    http.MethodPost,
			path:      "/update/gauge/test/invalid",
			response:  http.StatusBadRequest,
			setupMock: func(m *MockMetricStorage) {},
		},
		{
			name:     "UpdateMetrics returns StatusBadRequest for storage error",
			method:   http.MethodPost,
			path:     "/update/counter/test/123",
			response: http.StatusBadRequest,
			setupMock: func(m *MockMetricStorage) {
				m.On("UpdateMetric", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("domain.Metric")).
					Return(errors.New("storage error"))
			},
		},
		{
			name:     "GetAllMetrics returns StatusOK with no metrics",
			method:   http.MethodGet,
			path:     "/",
			response: http.StatusOK,
			setupMock: func(m *MockMetricStorage) {
				m.On("GetAllMetrics", mock.AnythingOfType("context.backgroundCtx")).
					Return([]domain.Metric{})
			},
		},
		{
			name:     "GetAllMetrics returns StatusOK with metrics",
			method:   http.MethodGet,
			path:     "/",
			response: http.StatusOK,
			setupMock: func(m *MockMetricStorage) {
				m.On("GetAllMetrics", mock.AnythingOfType("context.backgroundCtx")).
					Return([]domain.Metric{
						{
							ID:    "test_counter",
							MType: "counter",
							Delta: func() *int64 { v := int64(42); return &v }(),
						},
					})
			},
		},
		{
			name:      "Returns StatusNotFound for wrong method on update",
			method:    http.MethodGet,
			path:      "/update/gauge/test/123",
			response:  http.StatusNotFound,
			setupMock: func(m *MockMetricStorage) {},
		},
		{
			name:      "Returns StatusNotFound for wrong method on value",
			method:    http.MethodPost,
			path:      "/value/gauge/test",
			response:  http.StatusNotFound,
			setupMock: func(m *MockMetricStorage) {},
		},
		{
			name:      "Returns StatusNotFound for wrong path",
			method:    http.MethodPost,
			path:      "/record/gauge/test/123",
			response:  http.StatusNotFound,
			setupMock: func(m *MockMetricStorage) {},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockStorage := &MockMetricStorage{}
			tc.setupMock(mockStorage)

			handler := NewMetricsHandler(mockStorage)
			router := gin.New()

			// Register routes
			router.GET("/", handler.GetAllMetrics)
			router.GET("/value/:metricType/:metricName", handler.GetMetric)
			router.POST("/update/:metricType/:metricName/:metricValue", handler.UpdateMetrics)

			req, _ := http.NewRequest(tc.method, tc.path, nil)
			w := httptest.NewRecorder()

			// Act
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.response, w.Code)
			mockStorage.AssertExpectations(t)
		})
	}

	t.Run("constructor", func(t *testing.T) {
		mockStorage := &MockMetricStorage{}
		handler := NewMetricsHandler(mockStorage)

		assert.NotNil(t, handler)
		assert.Equal(t, mockStorage, handler.storage)
	})

	t.Run("interface compliance", func(t *testing.T) {
		mockStorage := &MockMetricStorage{}
		handler := NewMetricsHandler(mockStorage)

		assert.Implements(t, (*interface {
			GetMetric(*gin.Context)
			GetAllMetrics(*gin.Context)
			UpdateMetrics(*gin.Context)
		})(nil), handler)
	})
}
