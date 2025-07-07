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
