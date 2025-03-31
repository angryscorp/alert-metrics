package domain

type MetricMonitor interface {
	Start()
	Stop()
	GetCounters() map[string]int64
	GetGauges() map[string]float64
}
