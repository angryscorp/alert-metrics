package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/angryscorp/alert-metrics/internal/domain"
)

func TestMetricsJSONHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name        string
		method      string
		path        string
		body        interface{}
		contentType string
		response    int
		setupMock   func(*MockMetricStorage)
	}{
		{
			name:        "FetchMetricsJSON returns StatusOK for existing metric",
			method:      http.MethodPost,
			path:        "/value/",
			body:        domain.Metric{ID: "test_counter", MType: "counter"},
			contentType: "application/json",
			response:    http.StatusOK,
			setupMock: func(m *MockMetricStorage) {
				m.On("GetMetric", mock.AnythingOfType("context.backgroundCtx"), domain.MetricType("counter"), "test_counter").
					Return(domain.Metric{
						ID:    "test_counter",
						MType: "counter",
						Delta: func() *int64 { v := int64(42); return &v }(),
					}, true)
			},
		},
		{
			name:        "FetchMetricsJSON returns StatusNotFound for non-existing metric",
			method:      http.MethodPost,
			path:        "/value/",
			body:        domain.Metric{ID: "unknown", MType: "counter"},
			contentType: "application/json",
			response:    http.StatusNotFound,
			setupMock: func(m *MockMetricStorage) {
				m.On("GetMetric", mock.AnythingOfType("context.backgroundCtx"), domain.MetricType("counter"), "unknown").
					Return(domain.Metric{}, false)
			},
		},
		{
			name:        "FetchMetricsJSON returns StatusUnsupportedMediaType for wrong content type",
			method:      http.MethodPost,
			path:        "/value/",
			body:        domain.Metric{ID: "test", MType: "counter"},
			contentType: "text/plain",
			response:    http.StatusUnsupportedMediaType,
			setupMock:   func(m *MockMetricStorage) {},
		},
		{
			name:        "FetchMetricsJSON returns StatusBadRequest for invalid JSON",
			method:      http.MethodPost,
			path:        "/value/",
			body:        "invalid json",
			contentType: "application/json",
			response:    http.StatusBadRequest,
			setupMock:   func(m *MockMetricStorage) {},
		},
		{
			name:        "UpdateMetricsJSON returns StatusOK for valid metric",
			method:      http.MethodPost,
			path:        "/update/",
			body:        domain.Metric{ID: "test_counter", MType: "counter", Delta: func() *int64 { v := int64(42); return &v }()},
			contentType: "application/json",
			response:    http.StatusOK,
			setupMock: func(m *MockMetricStorage) {
				m.On("UpdateMetric", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("domain.Metric")).
					Return(nil)
				m.On("GetMetric", mock.Anything, domain.MetricType("counter"), "test_counter").
					Return(domain.Metric{
						ID:    "test_counter",
						MType: "counter",
						Delta: func() *int64 { v := int64(42); return &v }(),
					}, true)
			},
		},
		{
			name:        "UpdateMetricsJSON returns StatusUnsupportedMediaType for wrong content type",
			method:      http.MethodPost,
			path:        "/update/",
			body:        domain.Metric{ID: "test", MType: "counter"},
			contentType: "text/plain",
			response:    http.StatusUnsupportedMediaType,
			setupMock:   func(m *MockMetricStorage) {},
		},
		{
			name:        "UpdateMetricsJSON returns StatusBadRequest for invalid JSON",
			method:      http.MethodPost,
			path:        "/update/",
			body:        "invalid json",
			contentType: "application/json",
			response:    http.StatusBadRequest,
			setupMock:   func(m *MockMetricStorage) {},
		},
		{
			name:        "UpdateMetricsJSON returns StatusBadRequest for storage error",
			method:      http.MethodPost,
			path:        "/update/",
			body:        domain.Metric{ID: "test_counter", MType: "counter", Delta: func() *int64 { v := int64(42); return &v }()},
			contentType: "application/json",
			response:    http.StatusBadRequest,
			setupMock: func(m *MockMetricStorage) {
				m.On("UpdateMetric", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("domain.Metric")).
					Return(errors.New("storage error"))
			},
		},
		{
			name:   "BatchUpdateFetchMetrics returns StatusOK for valid metrics",
			method: http.MethodPost,
			path:   "/updates/",
			body: []domain.Metric{
				{ID: "test_counter", MType: "counter", Delta: func() *int64 { v := int64(42); return &v }()},
				{ID: "test_gauge", MType: "gauge", Value: func() *float64 { v := 3.14; return &v }()},
			},
			contentType: "application/json",
			response:    http.StatusOK,
			setupMock: func(m *MockMetricStorage) {
				m.On("UpdateMetrics", mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("[]domain.Metric")).
					Return(nil)
			},
		},
		{
			name:        "BatchUpdateFetchMetrics returns StatusUnsupportedMediaType for wrong content type",
			method:      http.MethodPost,
			path:        "/updates/",
			body:        []domain.Metric{},
			contentType: "text/plain",
			response:    http.StatusUnsupportedMediaType,
			setupMock:   func(m *MockMetricStorage) {},
		},
		{
			name:        "BatchUpdateFetchMetrics returns StatusBadRequest for invalid JSON",
			method:      http.MethodPost,
			path:        "/updates/",
			body:        "invalid json",
			contentType: "application/json",
			response:    http.StatusBadRequest,
			setupMock:   func(m *MockMetricStorage) {},
		},
		{
			name:        "BatchUpdateFetchMetrics returns StatusInternalServerError for storage error",
			method:      http.MethodPost,
			path:        "/updates/",
			body:        []domain.Metric{{ID: "test", MType: "counter"}},
			contentType: "application/json",
			response:    http.StatusInternalServerError,
			setupMock: func(m *MockMetricStorage) {
				m.On("UpdateMetrics", mock.Anything, mock.AnythingOfType("[]domain.Metric")).
					Return(errors.New("storage error"))
			},
		},
		{
			name:        "Returns StatusNotFound for wrong method",
			method:      http.MethodGet,
			path:        "/update/",
			body:        nil,
			contentType: "application/json",
			response:    http.StatusNotFound,
			setupMock:   func(m *MockMetricStorage) {},
		},
		{
			name:        "Returns StatusNotFound for wrong path",
			method:      http.MethodPost,
			path:        "/invalid/",
			body:        nil,
			contentType: "application/json",
			response:    http.StatusNotFound,
			setupMock:   func(m *MockMetricStorage) {},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockStorage := &MockMetricStorage{}
			tc.setupMock(mockStorage)

			handler := NewMetricsJSONHandler(mockStorage)
			router := gin.New()

			// Register routes
			router.POST("/value/", handler.FetchMetricsJSON)
			router.POST("/update/", handler.UpdateMetricsJSON)
			router.POST("/updates/", handler.BatchUpdateFetchMetrics)

			var body []byte
			if tc.body != nil {
				body, _ = json.Marshal(tc.body)
			}

			req, _ := http.NewRequest(tc.method, tc.path, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", tc.contentType)
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
		handler := NewMetricsJSONHandler(mockStorage)

		assert.NotNil(t, handler)
		assert.Equal(t, mockStorage, handler.storage)
	})

	t.Run("interface compliance", func(t *testing.T) {
		mockStorage := &MockMetricStorage{}
		handler := NewMetricsJSONHandler(mockStorage)

		assert.Implements(t, (*interface {
			FetchMetricsJSON(*gin.Context)
			UpdateMetricsJSON(*gin.Context)
			BatchUpdateFetchMetrics(*gin.Context)
		})(nil), handler)
	})
}
