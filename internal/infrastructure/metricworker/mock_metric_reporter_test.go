package metricworker

import (
	"github.com/stretchr/testify/mock"

	"github.com/angryscorp/alert-metrics/internal/domain"
)

type MockMetricReporter struct {
	mock.Mock
}

func (m *MockMetricReporter) ReportMetric(metric domain.Metric) {
	m.Called(metric)
}

func (m *MockMetricReporter) ReportBatch(metrics []domain.Metric) {
	m.Called(metrics)
}

func (m *MockMetricReporter) ReportRawMetric(metricType domain.MetricType, key string, value string) {
	m.Called(metricType, key, value)
}
