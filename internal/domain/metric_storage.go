package domain

import "context"

type MetricStorage interface {
	GetAllMetrics(ctx context.Context) []Metric
	UpdateMetric(ctx context.Context, metric Metric) error
	UpdateMetrics(ctx context.Context, metrics []Metric) error
	GetMetric(ctx context.Context, metricType MetricType, metricName string) (Metric, bool)
	Ping(ctx context.Context) error
}
