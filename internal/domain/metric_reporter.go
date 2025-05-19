package domain

type MetricReporter interface {
	ReportRawMetric(metricType MetricType, key string, value string)
	ReportMetric(metric Metric)
	ReportBatch(metrics []Metric)
}
