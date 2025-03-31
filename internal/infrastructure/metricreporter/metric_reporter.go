package metricreporter

import (
	"github.com/angryscorp/alert-metrics/internal/domain"
	"log/slog"
	"net/http"
	"os"
)

type HTTPMetricReporter struct {
	baseURL string
	client  *http.Client
	logger  *slog.Logger
}

func NewHTTPMetricReporter(baseURL string, client *http.Client) *HTTPMetricReporter {
	return &HTTPMetricReporter{
		baseURL: baseURL,
		client:  client,
		logger:  slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
}

func (mr *HTTPMetricReporter) Report(metricType domain.MetricType, key string, value string) {
	mr.logger.Info("report metric request", "metric type", metricType, "metric name", key, "metric value", value)

	resp, err := mr.client.Post(mr.baseURL+"/update/"+metricType+"/"+key+"/"+value, "text/plain", nil)
	if err != nil {
		mr.logger.Error("failed to report metric", "metric type", metricType, "metric name", key, "metric value", value, "error", err)
		return
	}
	_ = resp.Body.Close()

	mr.logger.Info("report metric response", "metric type", metricType, "metric name", key, "metric value", value, "status", resp.Status, "status code", resp.StatusCode)
}
