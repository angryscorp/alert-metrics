package metricmonitor

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewRuntimeMonitor(t *testing.T) {
	tests := []struct {
		name         string
		pollInterval time.Duration
	}{
		{
			name:         "create with 1 second interval",
			pollInterval: time.Second,
		},
		{
			name:         "create with 100ms interval",
			pollInterval: 100 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			monitor := NewRuntimeMonitor(tt.pollInterval)

			assert.NotNil(t, monitor)
			assert.Equal(t, tt.pollInterval, monitor.pollInterval)
			assert.False(t, monitor.isStarted)
			assert.NotNil(t, monitor.counters)
			assert.NotNil(t, monitor.gauges)
			assert.Empty(t, monitor.counters)
			assert.Empty(t, monitor.gauges)
		})
	}
}

func TestRuntimeMonitor_StartStop(t *testing.T) {
	tests := []struct {
		name         string
		pollInterval time.Duration
		operation    string
	}{
		{
			name:         "start monitor",
			pollInterval: 10 * time.Millisecond,
			operation:    "start",
		},
		{
			name:         "stop monitor",
			pollInterval: 10 * time.Millisecond,
			operation:    "stop",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			monitor := NewRuntimeMonitor(tt.pollInterval)

			switch tt.operation {
			case "start":
				monitor.Start()
				assert.True(t, monitor.isStarted)

				// Give it some time to collect metrics
				time.Sleep(100 * time.Millisecond)

				metrics := monitor.GetMetrics()
				assert.NotEmpty(t, metrics.Counters)
				assert.NotEmpty(t, metrics.Gauges)

				monitor.Stop()

			case "stop":
				monitor.Start()
				assert.True(t, monitor.isStarted)

				monitor.Stop()
				assert.False(t, monitor.isStarted)
			}
		})
	}
}

func TestRuntimeMonitor_GetMetrics(t *testing.T) {
	tests := []struct {
		name           string
		pollInterval   time.Duration
		startMonitor   bool
		waitTime       time.Duration
		expectedFields []string
	}{
		{
			name:           "get metrics from stopped monitor",
			pollInterval:   10 * time.Millisecond,
			startMonitor:   false,
			waitTime:       0,
			expectedFields: []string{},
		},
		{
			name:         "get metrics from running monitor",
			pollInterval: 10 * time.Millisecond,
			startMonitor: true,
			waitTime:     50 * time.Millisecond,
			expectedFields: []string{
				"PollCount", "RandomValue", "Alloc", "HeapAlloc", "Sys",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			monitor := NewRuntimeMonitor(tt.pollInterval)

			if tt.startMonitor {
				monitor.Start()
				defer monitor.Stop()
			}

			if tt.waitTime > 0 {
				time.Sleep(tt.waitTime)
			}

			metrics := monitor.GetMetrics()

			if len(tt.expectedFields) == 0 {
				assert.Empty(t, metrics.Counters)
				assert.Empty(t, metrics.Gauges)
			} else {
				assert.NotEmpty(t, metrics.Counters)
				assert.NotEmpty(t, metrics.Gauges)

				for _, field := range tt.expectedFields {
					if field == "PollCount" {
						assert.Contains(t, metrics.Counters, field)
					} else {
						assert.Contains(t, metrics.Gauges, field)
					}
				}
			}
		})
	}
}
