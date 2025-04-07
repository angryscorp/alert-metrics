package domain

type Metrics struct {
	Counters map[string]int64
	Gauges   map[string]float64
}

type MetricMonitor interface {
	Start()
	Stop()
	GetMetrics() Metrics
}
