package metricworker

import (
	"github.com/stretchr/testify/mock"

	"github.com/angryscorp/alert-metrics/internal/domain"
)

type MockMetricMonitor struct {
	mock.Mock
}

func (m *MockMetricMonitor) GetMetrics() domain.MetricsRawData {
	args := m.Called()
	return args.Get(0).(domain.MetricsRawData)
}

func (m *MockMetricMonitor) Start() {
	m.Called()
}

func (m *MockMetricMonitor) Stop() {
	m.Called()
}
