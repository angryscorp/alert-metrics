package domain

type MetricStorage interface {
	Get(metricType MetricType, key string) (string, bool)
	GetAllMetrics() map[string]string
	Update(metricType MetricType, key string, value string) error

	// TODO: Move to a new separate interface MetricJSONStorage
	UpdateMetrics(metrics Metrics) error
	GetMetrics(metricType MetricType, metricName string) Metrics
}
