package domain

type MetricStorage interface {
	Get(metricType MetricType, key string) (string, bool)
	GetAllMetrics() map[string]string
	Update(metricType MetricType, key string, value string) error
}
