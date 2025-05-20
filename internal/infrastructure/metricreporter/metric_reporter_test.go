package metricreporter

import (
	"github.com/angryscorp/alert-metrics/internal/domain"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type MockRoundTripper struct {
	lastRequest *http.Request
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	m.lastRequest = req
	return &http.Response{StatusCode: http.StatusOK}, nil
}

func Test_HTTPMetricReporter_Report(t *testing.T) {
	// Arrange
	transport := &MockRoundTripper{}
	mockClient := &http.Client{Transport: transport}
	reporter := NewHTTPMetricReporter("http://example.com", mockClient)

	// Act
	reporter.ReportRawMetric(domain.MetricTypeGauge, "key", "value")

	// Assert
	assert.Equal(t, transport.lastRequest.URL.String(), "http://example.com/update/gauge/key/value")
}
