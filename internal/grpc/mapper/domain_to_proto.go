package mapper

import (
	"github.com/angryscorp/alert-metrics/internal/domain"
	grpcmetrics "github.com/angryscorp/alert-metrics/internal/grpc/metrics"
)

func MetricToProto(metric domain.Metric) *grpcmetrics.Metric {
	protoMetric := &grpcmetrics.Metric{
		Id:   metric.ID,
		Type: MetricTypeToProto(metric.MType),
	}

	if metric.Delta != nil {
		protoMetric.Delta = metric.Delta
	}

	if metric.Value != nil {
		protoMetric.Value = metric.Value
	}

	return protoMetric
}

func MetricTypeToProto(metricType domain.MetricType) grpcmetrics.MetricType {
	switch metricType {
	case domain.MetricTypeCounter:
		return grpcmetrics.MetricType_METRIC_TYPE_COUNTER
	case domain.MetricTypeGauge:
		return grpcmetrics.MetricType_METRIC_TYPE_GAUGE
	default:
		return grpcmetrics.MetricType_METRIC_TYPE_UNSPECIFIED
	}
}
