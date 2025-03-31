package metricrouter

import (
	"fmt"
	"github.com/angryscorp/alert-metrics/internal/domain"
	"github.com/angryscorp/alert-metrics/internal/infrastructure/metricstorage"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type want struct {
	code        int
	contentType string
}

type metric struct {
	metricType  domain.MetricType
	metricName  string
	metricValue string
}

func (m *metric) toString() string {
	if m == nil {
		return ""
	}
	return fmt.Sprintf("/%s/%s/%s", m.metricType, m.metricName, m.metricValue)
}

func TestMetricRouter(t *testing.T) {

	router := NewMetricRouter(http.NewServeMux(), metricstorage.NewMemStorage())

	tests := []struct {
		name        string
		method      string
		path        string
		contentType string
		metric      *metric
		want        want
	}{
		{
			name:   "Healthcheck returns StatusOK in case of no errors",
			method: http.MethodGet,
			path:   "/health",
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "Update returns StatusOK in case of no errors",
			method:      http.MethodPost,
			path:        "/update",
			metric:      &metric{domain.MetricTypeGauge, "test", "123"},
			contentType: "text/plain",
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "Update returns StatusBadRequest in case of wrong method",
			method:      http.MethodGet,
			path:        "/update",
			metric:      &metric{domain.MetricTypeGauge, "test", "123"},
			contentType: "text/plain",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "Update returns StatusBadRequest in case of wrong content-type",
			method:      http.MethodGet,
			path:        "/update",
			metric:      &metric{domain.MetricTypeGauge, "test", "123"},
			contentType: "application/json",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "Update returns StatusBadRequest in case of wrong path",
			method:      http.MethodGet,
			path:        "/record",
			metric:      &metric{domain.MetricTypeGauge, "test", "123"},
			contentType: "text/plain",
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "Update returns StatusBadRequest in case of wrong value",
			method:      http.MethodGet,
			path:        "/update",
			metric:      &metric{domain.MetricTypeCounter, "test", "1.23"},
			contentType: "text/plain",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			request := httptest.NewRequest(test.method, test.path+test.metric.toString(), nil)
			request.Header.Set("Content-Type", test.contentType)
			w := httptest.NewRecorder()

			// Act
			router.ServeHTTP(w, request)
			res := w.Result()
			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(res.Body)

			// Assert
			require.NotNil(t, res)
			require.Equal(t, test.want.code, res.StatusCode)
			require.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
