package metricreporter

import (
	"bytes"
	"encoding/json"
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

var _ domain.MetricReporter = (*HTTPMetricReporter)(nil)

func NewHTTPMetricReporter(baseURL string, client *http.Client) *HTTPMetricReporter {
	return &HTTPMetricReporter{
		baseURL: baseURL,
		client:  client,
		logger:  slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
}

func (mr *HTTPMetricReporter) Report(metricType domain.MetricType, key string, value string) {
	mr.logger.Info("report metric request", "metric type", metricType, "metric name", key, "metric value", value)

	resp, err := mr.client.Post(mr.baseURL+"/update/"+string(metricType)+"/"+key+"/"+value, "text/plain", nil)
	if err != nil {
		mr.logger.Error("failed to report metric", "metric type", metricType, "metric name", key, "metric value", value, "error", err)
		return
	}
	_ = resp.Body.Close()

	mr.logger.Info("report metric response", "metric type", metricType, "metric name", key, "metric value", value, "status", resp.Status, "status code", resp.StatusCode)
}

func (mr *HTTPMetricReporter) ReportMetrics(metrics domain.Metric) {
	mr.logger.Info("report metric request", "metrics", metrics)

	bodyBytes, err := json.Marshal(metrics)
	if err != nil {
		mr.logger.Error("failed to convert metrics to json", "metrics", metrics)
		return
	}

	req, err := http.NewRequest(http.MethodPost, mr.baseURL+"/update/", bytes.NewBuffer(bodyBytes))
	if err != nil {
		mr.logger.Error("failed to build post request", "json", bodyBytes, "error", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := mr.client.Do(req)
	if err != nil {
		mr.logger.Error("failed to report metrics", "metrics", metrics, "error", err)
		return
	}
	_ = resp.Body.Close()

	mr.logger.Info("report metric response", "metrics", metrics)
}
