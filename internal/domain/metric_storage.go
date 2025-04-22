package domain

type MetricStorage interface {
	GetAllMetrics() MetricRepresentatives
	UpdateMetrics(metrics Metrics) error
	GetMetrics(metricType MetricType, metricName string) (Metrics, bool)
}
