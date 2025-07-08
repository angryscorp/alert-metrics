package domain

// MetricReporter defines an interface for reporting metrics including individual, raw, or batch metric data.
type MetricReporter interface {
	ReportRawMetric(metricType MetricType, key string, value string)
	ReportMetric(metric Metric)
	ReportBatch(metrics []Metric)
}
