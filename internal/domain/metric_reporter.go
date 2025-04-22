package domain

type MetricReporter interface {
	Report(metricType MetricType, key string, value string)
	ReportMetrics(metrics Metric)
}
