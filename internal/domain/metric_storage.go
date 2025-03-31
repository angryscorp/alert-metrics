package domain

type MetricStorage interface {
	Get(metricType MetricType, key string) (string, error)
	Update(metricType MetricType, key string, value string) error
}
