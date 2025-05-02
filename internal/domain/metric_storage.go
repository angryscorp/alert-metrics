package domain

type MetricStorage interface {
	GetAllMetrics() []Metric
	UpdateMetric(metric Metric) error
	GetMetric(metricType MetricType, metricName string) (Metric, bool)
	Ping() error
}
