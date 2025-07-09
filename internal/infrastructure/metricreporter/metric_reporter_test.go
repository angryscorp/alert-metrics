package metricreporter

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/angryscorp/alert-metrics/internal/domain"
)

type MockRoundTripper struct {
	lastRequest *http.Request
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	m.lastRequest = req
	return &http.Response{StatusCode: http.StatusOK}, nil
}

func Test_HTTPMetricReporter_ReportRawMetric(t *testing.T) {
	// Arrange
	transport := &MockRoundTripper{}
	mockClient := &http.Client{Transport: transport}
	reporter := NewHTTPMetricReporter("http://example.com", mockClient)

	// Act
	reporter.ReportRawMetric(domain.MetricTypeGauge, "key", "value")

	// Assert
	assert.Equal(t, transport.lastRequest.URL.String(), "http://example.com/update/gauge/key/value")
}

func Test_HTTPMetricReporter_ReportMetric(t *testing.T) {
	// Arrange
	transport := &MockRoundTripper{}
	mockClient := &http.Client{Transport: transport}
	reporter := NewHTTPMetricReporter("http://example.com", mockClient)

	gaugeValue := 42.5
	metric := domain.Metric{
		ID:    "test_gauge",
		MType: domain.MetricTypeGauge,
		Value: &gaugeValue,
	}

	// Act
	reporter.ReportMetric(metric)

	// Assert
	assert.Equal(t, "http://example.com/update/", transport.lastRequest.URL.String())
	assert.Equal(t, "POST", transport.lastRequest.Method)
	assert.Equal(t, "application/json", transport.lastRequest.Header.Get("Content-Type"))
}

func Test_HTTPMetricReporter_ReportBatch(t *testing.T) {
	// Arrange
	transport := &MockRoundTripper{}
	mockClient := &http.Client{Transport: transport}
	reporter := NewHTTPMetricReporter("http://example.com", mockClient)

	gaugeValue := 42.5
	counterValue := int64(10)
	metrics := []domain.Metric{
		{
			ID:    "test_gauge",
			MType: domain.MetricTypeGauge,
			Value: &gaugeValue,
		},
		{
			ID:    "test_counter",
			MType: domain.MetricTypeCounter,
			Delta: &counterValue,
		},
	}

	// Act
	reporter.ReportBatch(metrics)

	// Assert
	assert.Equal(t, "http://example.com/updates/", transport.lastRequest.URL.String())
	assert.Equal(t, "POST", transport.lastRequest.Method)
	assert.Equal(t, "application/json", transport.lastRequest.Header.Get("Content-Type"))
}
