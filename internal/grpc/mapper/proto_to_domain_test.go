package mapper

import (
	"testing"

	"github.com/angryscorp/alert-metrics/internal/domain"
	grpcmetrics "github.com/angryscorp/alert-metrics/internal/grpc/metrics"
)

func TestMetricToDomain(t *testing.T) {
	tests := []struct {
		name  string
		input *grpcmetrics.Metric
		want  domain.Metric
	}{
		{
			name: "counter metric",
			input: &grpcmetrics.Metric{
				Id:    "test_counter",
				Type:  grpcmetrics.MetricType_METRIC_TYPE_COUNTER,
				Delta: &[]int64{42}[0],
			},
			want: domain.Metric{
				ID:    "test_counter",
				MType: domain.MetricTypeCounter,
				Delta: &[]int64{42}[0],
			},
		},
		{
			name: "gauge metric",
			input: &grpcmetrics.Metric{
				Id:    "test_gauge",
				Type:  grpcmetrics.MetricType_METRIC_TYPE_GAUGE,
				Value: &[]float64{3.14}[0],
			},
			want: domain.Metric{
				ID:    "test_gauge",
				MType: domain.MetricTypeGauge,
				Value: &[]float64{3.14}[0],
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MetricToDomain(tt.input)

			if got.ID != tt.want.ID {
				t.Errorf("MetricToDomain().ID = %v, want %v", got.ID, tt.want.ID)
			}

			if got.MType != tt.want.MType {
				t.Errorf("MetricToDomain().MType = %v, want %v", got.MType, tt.want.MType)
			}

			if !equalInt64Ptr(got.Delta, tt.want.Delta) {
				t.Errorf("MetricToDomain().Delta = %v, want %v", ptrValue(got.Delta), ptrValue(tt.want.Delta))
			}

			if !equalFloat64Ptr(got.Value, tt.want.Value) {
				t.Errorf("MetricToDomain().Value = %v, want %v", ptrValue(got.Value), ptrValue(tt.want.Value))
			}
		})
	}
}

func TestMetricTypeToDomain(t *testing.T) {
	tests := []struct {
		name  string
		input grpcmetrics.MetricType
		want  domain.MetricType
	}{
		{
			name:  "counter",
			input: grpcmetrics.MetricType_METRIC_TYPE_COUNTER,
			want:  domain.MetricTypeCounter,
		},
		{
			name:  "gauge",
			input: grpcmetrics.MetricType_METRIC_TYPE_GAUGE,
			want:  domain.MetricTypeGauge,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MetricTypeToDomain(tt.input)
			if got != tt.want {
				t.Errorf("MetricTypeToDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}
