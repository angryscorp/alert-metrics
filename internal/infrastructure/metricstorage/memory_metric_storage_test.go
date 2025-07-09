package metricstorage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/angryscorp/alert-metrics/internal/domain"
)

func TestMemoryMetricStorage_UpdateMetric(t *testing.T) {
	ctx := context.Background()
	counterValue := int64(10)
	gaugeValue := float64(3.14)

	tests := []struct {
		name        string
		metric      domain.Metric
		expectError bool
		errorMsg    string
	}{
		{
			name: "update counter metric",
			metric: domain.Metric{
				ID:    "test_counter",
				MType: domain.MetricTypeCounter,
				Delta: &counterValue,
			},
			expectError: false,
		},
		{
			name: "update gauge metric",
			metric: domain.Metric{
				ID:    "test_gauge",
				MType: domain.MetricTypeGauge,
				Value: &gaugeValue,
			},
			expectError: false,
		},
		{
			name: "unsupported metric type",
			metric: domain.Metric{
				ID:    "test_invalid",
				MType: "invalid",
			},
			expectError: true,
			errorMsg:    "unsupported metric type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := NewMemoryMetricStorage()

			err := storage.UpdateMetric(ctx, tt.metric)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMemoryMetricStorage_GetMetric(t *testing.T) {
	ctx := context.Background()
	counterValue := int64(10)
	gaugeValue := float64(3.14)

	tests := []struct {
		name        string
		setupData   []domain.Metric
		metricType  domain.MetricType
		metricName  string
		expectFound bool
		expected    domain.Metric
	}{
		{
			name: "get existing counter metric",
			setupData: []domain.Metric{
				{
					ID:    "test_counter",
					MType: domain.MetricTypeCounter,
					Delta: &counterValue,
				},
			},
			metricType:  domain.MetricTypeCounter,
			metricName:  "test_counter",
			expectFound: true,
			expected: domain.Metric{
				ID:    "test_counter",
				MType: domain.MetricTypeCounter,
				Delta: &counterValue,
			},
		},
		{
			name: "get existing gauge metric",
			setupData: []domain.Metric{
				{
					ID:    "test_gauge",
					MType: domain.MetricTypeGauge,
					Value: &gaugeValue,
				},
			},
			metricType:  domain.MetricTypeGauge,
			metricName:  "test_gauge",
			expectFound: true,
			expected: domain.Metric{
				ID:    "test_gauge",
				MType: domain.MetricTypeGauge,
				Value: &gaugeValue,
			},
		},
		{
			name:        "get non-existing metric",
			setupData:   []domain.Metric{},
			metricType:  domain.MetricTypeCounter,
			metricName:  "non_existing",
			expectFound: false,
			expected: domain.Metric{
				ID:    "non_existing",
				MType: domain.MetricTypeCounter,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := NewMemoryMetricStorage()

			for _, metric := range tt.setupData {
				err := storage.UpdateMetric(ctx, metric)
				require.NoError(t, err)
			}

			result, found := storage.GetMetric(ctx, tt.metricType, tt.metricName)

			assert.Equal(t, tt.expectFound, found)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMemoryMetricStorage_GetAllMetrics(t *testing.T) {
	ctx := context.Background()
	storage := NewMemoryMetricStorage()

	counterValue := int64(10)
	gaugeValue := float64(3.14)

	metrics := []domain.Metric{
		{
			ID:    "counter1",
			MType: domain.MetricTypeCounter,
			Delta: &counterValue,
		},
		{
			ID:    "gauge1",
			MType: domain.MetricTypeGauge,
			Value: &gaugeValue,
		},
	}

	for _, metric := range metrics {
		err := storage.UpdateMetric(ctx, metric)
		require.NoError(t, err)
	}

	result := storage.GetAllMetrics(ctx)
	assert.Len(t, result, 2)
	assert.ElementsMatch(t, metrics, result)
}

func TestMemoryMetricStorage_UpdateMetrics(t *testing.T) {
	ctx := context.Background()
	counterValue := int64(10)
	gaugeValue := float64(3.14)

	tests := []struct {
		name        string
		metrics     []domain.Metric
		expectError bool
		errorMsg    string
	}{
		{
			name: "update multiple valid metrics",
			metrics: []domain.Metric{
				{
					ID:    "counter1",
					MType: domain.MetricTypeCounter,
					Delta: &counterValue,
				},
				{
					ID:    "gauge1",
					MType: domain.MetricTypeGauge,
					Value: &gaugeValue,
				},
			},
			expectError: false,
		},
		{
			name:        "empty metrics slice",
			metrics:     []domain.Metric{},
			expectError: false,
		},
		{
			name: "invalid metric in batch",
			metrics: []domain.Metric{
				{
					ID:    "invalid",
					MType: "invalid_type",
				},
			},
			expectError: true,
			errorMsg:    "unsupported metric type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := NewMemoryMetricStorage()

			err := storage.UpdateMetrics(ctx, tt.metrics)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMemoryMetricStorage_CounterAccumulation(t *testing.T) {
	ctx := context.Background()
	storage := NewMemoryMetricStorage()

	delta1 := int64(5)
	delta2 := int64(10)

	// First update
	metric := domain.Metric{
		ID:    "test_counter",
		MType: domain.MetricTypeCounter,
		Delta: &delta1,
	}
	err := storage.UpdateMetric(ctx, metric)
	require.NoError(t, err)

	result, found := storage.GetMetric(ctx, domain.MetricTypeCounter, "test_counter")
	require.True(t, found)
	assert.Equal(t, int64(5), *result.Delta)

	// Second update - should accumulate
	metric.Delta = &delta2
	err = storage.UpdateMetric(ctx, metric)
	require.NoError(t, err)

	result, found = storage.GetMetric(ctx, domain.MetricTypeCounter, "test_counter")
	require.True(t, found)
	assert.Equal(t, int64(15), *result.Delta)
}

func TestMemoryMetricStorage_GaugeReplacement(t *testing.T) {
	ctx := context.Background()
	storage := NewMemoryMetricStorage()

	value1 := float64(3.14)
	value2 := float64(2.71)

	// First update
	metric := domain.Metric{
		ID:    "test_gauge",
		MType: domain.MetricTypeGauge,
		Value: &value1,
	}
	err := storage.UpdateMetric(ctx, metric)
	require.NoError(t, err)

	result, found := storage.GetMetric(ctx, domain.MetricTypeGauge, "test_gauge")
	require.True(t, found)
	assert.Equal(t, 3.14, *result.Value)

	// Second update - should replace
	metric.Value = &value2
	err = storage.UpdateMetric(ctx, metric)
	require.NoError(t, err)

	result, found = storage.GetMetric(ctx, domain.MetricTypeGauge, "test_gauge")
	require.True(t, found)
	assert.Equal(t, 2.71, *result.Value)
}

func TestMemoryMetricStorage_Ping(t *testing.T) {
	ctx := context.Background()
	storage := NewMemoryMetricStorage()

	err := storage.Ping(ctx)
	assert.NoError(t, err)
}
