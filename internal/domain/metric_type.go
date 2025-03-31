package domain

type MetricType = string

var MetricTypes = []MetricType{
	MetricTypeCounter,
	MetricTypeGauge,
}

const (
	MetricTypeCounter MetricType = "counter" // new value increases the existing one
	MetricTypeGauge   MetricType = "gauge"   // new value replaces the previous one
)
