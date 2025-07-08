package domain

import "fmt"

// MetricType represents the type of a metric, such as "counter" or "gauge".
type MetricType string

// MetricTypeCounter represents a metric type where values are cumulative and increase over time.
// MetricTypeGauge represents a metric type where values can replace the previous ones.
const (
	MetricTypeCounter MetricType = "counter" // new value increases the existing one
	MetricTypeGauge   MetricType = "gauge"   // new value replaces the previous one
)

// MetricTypes is a list of predefined MetricType constants representing supported metric types: counter and gauge.
var MetricTypes = []MetricType{
	MetricTypeCounter,
	MetricTypeGauge,
}

// NewMetricType converts a string to a MetricType if it matches a valid predefined type, returning an error for invalid input.
func NewMetricType(s string) (MetricType, error) {
	for _, t := range MetricTypes {
		if string(t) == s {
			return t, nil
		}
	}
	return "", fmt.Errorf("invalid metric type")
}
