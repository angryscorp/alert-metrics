package mapper

import (
	"testing"

	"github.com/angryscorp/alert-metrics/internal/domain"
	grpcmetrics "github.com/angryscorp/alert-metrics/internal/grpc/metrics"
)

func TestMetricToProto(t *testing.T) {
	tests := []struct {
		name  string
		input domain.Metric
		want  *grpcmetrics.Metric
	}{
		{
			name: "counter metric",
			input: domain.Metric{
				ID:    "test_counter",
				MType: domain.MetricTypeCounter,
				Delta: &[]int64{42}[0],
			},
			want: &grpcmetrics.Metric{
				Id:    "test_counter",
				Type:  grpcmetrics.MetricType_METRIC_TYPE_COUNTER,
				Delta: &[]int64{42}[0],
			},
		},
		{
			name: "gauge metric",
			input: domain.Metric{
				ID:    "test_gauge",
				MType: domain.MetricTypeGauge,
				Value: &[]float64{3.14}[0],
			},
			want: &grpcmetrics.Metric{
				Id:    "test_gauge",
				Type:  grpcmetrics.MetricType_METRIC_TYPE_GAUGE,
				Value: &[]float64{3.14}[0],
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MetricToProto(tt.input)

			if got.Id != tt.want.Id {
				t.Errorf("MetricToProto().Id = %v, want %v", got.Id, tt.want.Id)
			}

			if got.Type != tt.want.Type {
				t.Errorf("MetricToProto().Type = %v, want %v", got.Type, tt.want.Type)
			}

			if !equalInt64Ptr(got.Delta, tt.want.Delta) {
				t.Errorf("MetricToProto().Delta = %v, want %v", ptrValue(got.Delta), ptrValue(tt.want.Delta))
			}

			if !equalFloat64Ptr(got.Value, tt.want.Value) {
				t.Errorf("MetricToProto().Value = %v, want %v", ptrValue(got.Value), ptrValue(tt.want.Value))
			}
		})
	}
}

func TestMetricTypeToProto(t *testing.T) {
	tests := []struct {
		name  string
		input domain.MetricType
		want  grpcmetrics.MetricType
	}{
		{
			name:  "counter",
			input: domain.MetricTypeCounter,
			want:  grpcmetrics.MetricType_METRIC_TYPE_COUNTER,
		},
		{
			name:  "gauge",
			input: domain.MetricTypeGauge,
			want:  grpcmetrics.MetricType_METRIC_TYPE_GAUGE,
		},
		{
			name:  "unknown",
			input: domain.MetricType("unknown"),
			want:  grpcmetrics.MetricType_METRIC_TYPE_UNSPECIFIED,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MetricTypeToProto(tt.input)
			if got != tt.want {
				t.Errorf("MetricTypeToProto() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper functions
func equalInt64Ptr(a, b *int64) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func equalFloat64Ptr(a, b *float64) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func ptrValue(ptr interface{}) interface{} {
	if ptr == nil {
		return nil
	}
	switch p := ptr.(type) {
	case *int64:
		return *p
	case *float64:
		return *p
	}
	return ptr
}
