package domain

import "context"

// MetricStorage defines an interface for managing and interacting with metrics in storage.
// GetAllMetrics retrieves all stored metrics.
// UpdateMetric updates a single metric in storage.
// UpdateMetrics updates multiple metrics in storage.
// GetMetric retrieves a specific metric by type and name. Returns the metric and a boolean indicating if found.
// Ping checks the liveness of the storage connection.
type MetricStorage interface {
	GetAllMetrics(ctx context.Context) []Metric
	UpdateMetric(ctx context.Context, metric Metric) error
	UpdateMetrics(ctx context.Context, metrics []Metric) error
	GetMetric(ctx context.Context, metricType MetricType, metricName string) (Metric, bool)
	Ping(ctx context.Context) error
}
