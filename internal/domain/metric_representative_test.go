package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricRepresentative_String(t *testing.T) {
	tests := []struct {
		name     string
		metric   MetricRepresentative
		expected string
	}{
		{
			name: "gauge metric",
			metric: MetricRepresentative{
				Type:  MetricTypeGauge,
				Name:  "cpu_usage",
				Value: "75.5",
			},
			expected: "cpu_usage (gauge) = 75.5",
		},
		{
			name: "counter metric",
			metric: MetricRepresentative{
				Type:  MetricTypeCounter,
				Name:  "requests_total",
				Value: "1000",
			},
			expected: "requests_total (counter) = 1000",
		},
		{
			name: "empty values",
			metric: MetricRepresentative{
				Type:  MetricTypeGauge,
				Name:  "",
				Value: "",
			},
			expected: " (gauge) = ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.metric.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewMetricRepresentatives(t *testing.T) {
	tests := []struct {
		name     string
		metrics  []Metric
		expected MetricRepresentatives
	}{
		{
			name:     "empty slice",
			metrics:  []Metric{},
			expected: MetricRepresentatives{},
		},
		{
			name: "mixed metrics",
			metrics: []Metric{
				{
					ID:    "cpu_usage",
					MType: MetricTypeGauge,
					Value: func() *float64 { v := 75.5; return &v }(),
				},
				{
					ID:    "requests_total",
					MType: MetricTypeCounter,
					Delta: func() *int64 { v := int64(1000); return &v }(),
				},
			},
			expected: MetricRepresentatives{
				{
					Type:  MetricTypeGauge,
					Name:  "cpu_usage",
					Value: "75.5",
				},
				{
					Type:  MetricTypeCounter,
					Name:  "requests_total",
					Value: "1000",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewMetricRepresentatives(tt.metrics)
			assert.Equal(t, tt.expected, result)
			assert.Len(t, result, len(tt.metrics))
		})
	}
}

func TestMetricRepresentatives_SortByName(t *testing.T) {
	tests := []struct {
		name     string
		metrics  MetricRepresentatives
		expected MetricRepresentatives
	}{
		{
			name:     "empty slice",
			metrics:  MetricRepresentatives{},
			expected: MetricRepresentatives{},
		},
		{
			name: "unsorted metrics",
			metrics: MetricRepresentatives{
				{
					Type:  MetricTypeGauge,
					Name:  "memory_usage",
					Value: "85.2",
				},
				{
					Type:  MetricTypeCounter,
					Name:  "api_requests",
					Value: "1000",
				},
				{
					Type:  MetricTypeGauge,
					Name:  "cpu_usage",
					Value: "75.5",
				},
			},
			expected: MetricRepresentatives{
				{
					Type:  MetricTypeCounter,
					Name:  "api_requests",
					Value: "1000",
				},
				{
					Type:  MetricTypeGauge,
					Name:  "cpu_usage",
					Value: "75.5",
				},
				{
					Type:  MetricTypeGauge,
					Name:  "memory_usage",
					Value: "85.2",
				},
			},
		},
		{
			name: "identical names",
			metrics: MetricRepresentatives{
				{
					Type:  MetricTypeGauge,
					Name:  "same_name",
					Value: "1",
				},
				{
					Type:  MetricTypeCounter,
					Name:  "same_name",
					Value: "2",
				},
			},
			expected: MetricRepresentatives{
				{
					Type:  MetricTypeGauge,
					Name:  "same_name",
					Value: "1",
				},
				{
					Type:  MetricTypeCounter,
					Name:  "same_name",
					Value: "2",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalSlice := tt.metrics
			result := tt.metrics.SortByName()

			assert.Equal(t, tt.expected, result)
			assert.Equal(t, tt.expected, originalSlice)
		})
	}
}
