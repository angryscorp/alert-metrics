package domain

type MetricStorage interface {
	GetAllMetrics() MetricRepresentatives
	UpdateMetrics(metrics Metric) error
	GetMetrics(metricType MetricType, metricName string) (Metric, bool)
}
