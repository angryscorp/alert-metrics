package domain

type MetricStorage interface {
	GetAllMetrics() []MetricRepresentative
	UpdateMetrics(metrics Metrics) error
	GetMetrics(metricType MetricType, metricName string) (Metrics, bool)
}
