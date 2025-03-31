package domain

import "fmt"

type MetricType = string

const (
	MetricTypeCounter MetricType = "counter" // new value increases the existing one
	MetricTypeGauge   MetricType = "gauge"   // new value replaces the previous one
)

var MetricTypes = []MetricType{
	MetricTypeCounter,
	MetricTypeGauge,
}

func NewMetricType(s string) (MetricType, error) {
	for _, t := range MetricTypes {
		if t == s {
			return t, nil
		}
	}
	return "", fmt.Errorf("invalid metric type")
}
