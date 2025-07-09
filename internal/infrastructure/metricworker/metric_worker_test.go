package metricworker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/angryscorp/alert-metrics/internal/domain"
)

func TestNewMetricWorker(t *testing.T) {
	tests := []struct {
		name           string
		reportInterval time.Duration
		rateLimiter    int
	}{
		{
			name:           "create with 1 second interval and 5 workers",
			reportInterval: time.Second,
			rateLimiter:    5,
		},
		{
			name:           "create with 100ms interval and 1 worker",
			reportInterval: 100 * time.Millisecond,
			rateLimiter:    1,
		},
		{
			name:           "create with 10ms interval and 10 workers",
			reportInterval: 10 * time.Millisecond,
			rateLimiter:    10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockMonitor := &MockMetricMonitor{}
			mockReporter := &MockMetricReporter{}

			worker := NewMetricWorker(mockMonitor, mockReporter, tt.reportInterval, tt.rateLimiter)

			assert.NotNil(t, worker)
			assert.Equal(t, mockMonitor, worker.metricMonitor)
			assert.Equal(t, mockReporter, worker.metricReporter)
			assert.Equal(t, tt.reportInterval, worker.reportInterval)
			assert.Equal(t, tt.rateLimiter, worker.rateLimiter)
			assert.False(t, worker.isRunning)
			assert.NotNil(t, worker.requestChan)
		})
	}
}

func TestMetricWorker_StartStop(t *testing.T) {
	mockMonitor := &MockMetricMonitor{}
	mockReporter := &MockMetricReporter{}

	mockMonitor.
		On("GetMetrics").
		Return(domain.MetricsRawData{
			Counters: map[string]int64{"test_counter": 1},
			Gauges:   map[string]float64{"test_gauge": 2},
		})

	mockReporter.
		On("ReportBatch", mock.AnythingOfType("[]domain.Metric")).
		Return()

	worker := NewMetricWorker(mockMonitor, mockReporter, time.Second, 5)

	worker.Start()
	assert.True(t, worker.isRunning)

	worker.Stop()
	assert.False(t, worker.isRunning)
}
