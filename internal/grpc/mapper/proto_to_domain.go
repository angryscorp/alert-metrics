package mapper

import (
	"github.com/angryscorp/alert-metrics/internal/domain"
	grpcmetrics "github.com/angryscorp/alert-metrics/internal/grpc/metrics"
)

func MetricToDomain(protoMetric *grpcmetrics.Metric) domain.Metric {
	metric := domain.Metric{
		ID:    protoMetric.Id,
		MType: MetricTypeToDomain(protoMetric.Type),
	}

	if protoMetric.Delta != nil {
		metric.Delta = protoMetric.Delta
	}

	if protoMetric.Value != nil {
		metric.Value = protoMetric.Value
	}

	return metric
}

func MetricTypeToDomain(protoType grpcmetrics.MetricType) domain.MetricType {
	switch protoType {
	case grpcmetrics.MetricType_METRIC_TYPE_COUNTER:
		return domain.MetricTypeCounter
	case grpcmetrics.MetricType_METRIC_TYPE_GAUGE:
		return domain.MetricTypeGauge
	default:
		return domain.MetricTypeCounter // default fallback
	}
}
