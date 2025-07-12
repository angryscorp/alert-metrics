package domain

// MetricsRawData represents a structure to store raw metrics data, including counters and gauges.
// Counters are integer-based metrics representing accumulated counts.
// Gauges are floating-point metrics representing point-in-time values.
type MetricsRawData struct {
	Counters map[string]int64
	Gauges   map[string]float64
}

// MetricMonitor defines methods for starting, stopping, and retrieving metrics data from a monitoring system.
type MetricMonitor interface {
	Start()
	Stop()
	GetMetrics() MetricsRawData
}
