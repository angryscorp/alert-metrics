package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMetrics(t *testing.T) {
	tests := []struct {
		name          string
		metricType    string
		metricName    string
		value         string
		expectError   bool
		expectedID    string
		expectedType  MetricType
		expectedDelta *int64
		expectedValue *float64
	}{
		{
			name:          "valid counter metric",
			metricType:    "counter",
			metricName:    "test_counter",
			value:         "123",
			expectError:   false,
			expectedID:    "test_counter",
			expectedType:  MetricTypeCounter,
			expectedDelta: pointerFrom[int64](123),
		},
		{
			name:          "valid gauge metric",
			metricType:    "gauge",
			metricName:    "test_gauge",
			value:         "45.67",
			expectError:   false,
			expectedID:    "test_gauge",
			expectedType:  MetricTypeGauge,
			expectedValue: pointerFrom(45.67),
		},
		{
			name:          "counter with zero value",
			metricType:    "counter",
			metricName:    "zero_counter",
			value:         "0",
			expectError:   false,
			expectedID:    "zero_counter",
			expectedType:  MetricTypeCounter,
			expectedDelta: pointerFrom[int64](0),
		},
		{
			name:          "gauge with negative value",
			metricType:    "gauge",
			metricName:    "negative_gauge",
			value:         "-12.34",
			expectError:   false,
			expectedID:    "negative_gauge",
			expectedType:  MetricTypeGauge,
			expectedValue: pointerFrom(-12.34),
		},
		{
			name:        "invalid metric type",
			metricType:  "invalid",
			metricName:  "test_metric",
			value:       "123",
			expectError: true,
		},
		{
			name:        "empty metric name",
			metricType:  "counter",
			metricName:  "",
			value:       "123",
			expectError: true,
		},
		{
			name:        "empty value",
			metricType:  "counter",
			metricName:  "test_metric",
			value:       "",
			expectError: true,
		},
		{
			name:        "invalid counter value",
			metricType:  "counter",
			metricName:  "test_counter",
			value:       "not_a_number",
			expectError: true,
		},
		{
			name:        "invalid gauge value",
			metricType:  "gauge",
			metricName:  "test_gauge",
			value:       "not_a_number",
			expectError: true,
		},
		{
			name:          "gauge with scientific notation",
			metricType:    "gauge",
			metricName:    "scientific_gauge",
			value:         "1.23e-4",
			expectError:   false,
			expectedID:    "scientific_gauge",
			expectedType:  MetricTypeGauge,
			expectedValue: pointerFrom(1.23e-4),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metric, err := NewMetrics(tt.metricType, tt.metricName, tt.value)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, metric)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, metric)

			assert.Equal(t, tt.expectedID, metric.ID)
			assert.Equal(t, tt.expectedType, metric.MType)

			if tt.expectedDelta != nil {
				require.NotNil(t, metric.Delta)
				assert.Equal(t, *tt.expectedDelta, *metric.Delta)
				assert.Nil(t, metric.Value)
			}

			if tt.expectedValue != nil {
				require.NotNil(t, metric.Value)
				assert.Equal(t, *tt.expectedValue, *metric.Value)
				assert.Nil(t, metric.Delta)
			}
		})
	}
}

func TestMetric_StringValue(t *testing.T) {
	tests := []struct {
		name     string
		metric   Metric
		expected string
	}{
		{
			name: "gauge metric with float",
			metric: Metric{
				ID:    "test_gauge",
				MType: MetricTypeGauge,
				Value: pointerFrom(123.456),
			},
			expected: "123.456",
		},
		{
			name: "gauge metric with zero",
			metric: Metric{
				ID:    "zero_gauge",
				MType: MetricTypeGauge,
				Value: pointerFrom(0.0),
			},
			expected: "0",
		},
		{
			name: "gauge metric with negative value",
			metric: Metric{
				ID:    "negative_gauge",
				MType: MetricTypeGauge,
				Value: pointerFrom(-45.67),
			},
			expected: "-45.67",
		},
		{
			name: "counter metric with positive integer",
			metric: Metric{
				ID:    "test_counter",
				MType: MetricTypeCounter,
				Delta: pointerFrom[int64](789),
			},
			expected: "789",
		},
		{
			name: "counter metric with zero",
			metric: Metric{
				ID:    "zero_counter",
				MType: MetricTypeCounter,
				Delta: pointerFrom[int64](0),
			},
			expected: "0",
		},
		{
			name: "counter metric with negative integer",
			metric: Metric{
				ID:    "negative_counter",
				MType: MetricTypeCounter,
				Delta: pointerFrom[int64](-123),
			},
			expected: "-123",
		},
		{
			name: "invalid metric type",
			metric: Metric{
				ID:    "invalid_metric",
				MType: MetricType("invalid"),
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.metric.StringValue()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMetric_StringValue_NilPointers(t *testing.T) {
	t.Run("gauge metric with nil value panics", func(t *testing.T) {
		metric := Metric{
			ID:    "nil_gauge",
			MType: MetricTypeGauge,
			Value: nil,
		}

		assert.Panics(t, func() {
			metric.StringValue()
		})
	})

	t.Run("counter metric with nil delta panics", func(t *testing.T) {
		metric := Metric{
			ID:    "nil_counter",
			MType: MetricTypeCounter,
			Delta: nil,
		}

		assert.Panics(t, func() {
			metric.StringValue()
		})
	})
}

func pointerFrom[T any](v T) *T {
	return &v
}
