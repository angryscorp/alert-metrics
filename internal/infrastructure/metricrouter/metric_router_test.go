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
			name:   "Healthcheck returns 200 in case of no errors",
			method: http.MethodGet,
			path:   "/health",
			want: want{
				code:        200,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:        "Update returns 200 in case of no errors",
			method:      http.MethodPost,
			path:        "/update",
			metric:      &metric{domain.MetricTypeGauge, "test", "123"},
			contentType: "text/plain",
			want: want{
				code:        200,
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
