package domain

type MetricStorage interface {
	GetAllMetrics() map[string]string
	UpdateMetrics(metrics Metrics) error
	GetMetrics(metricType MetricType, metricName string) (Metrics, bool)
}
