package metricmonitor

import (
	"github.com/angryscorp/alert-metrics/internal/domain"
	"math/rand/v2"
	"runtime"
	"time"
)

type RuntimeMonitor struct {
	counters     map[string]int64
	gauges       map[string]float64
	pollInterval time.Duration
	isStarted    bool

	m runtime.MemStats
}

func NewRuntimeMonitor(pollInterval time.Duration) *RuntimeMonitor {
	return &RuntimeMonitor{
		counters:     make(map[string]int64),
		gauges:       make(map[string]float64),
		pollInterval: pollInterval,
	}
}

var _ domain.MetricMonitor = (*RuntimeMonitor)(nil)

func (m *RuntimeMonitor) update() {
	// Read the current values
	runtime.ReadMemStats(&m.m)

	// Runtime metrics
	m.gauges["Alloc"] = float64(m.m.Alloc)
	m.gauges["BuckHashSys"] = float64(m.m.BuckHashSys)
	m.gauges["Frees"] = float64(m.m.Frees)
	m.gauges["GCCPUFraction"] = m.m.GCCPUFraction
	m.gauges["GCSys"] = float64(m.m.GCSys)
	m.gauges["HeapAlloc"] = float64(m.m.HeapAlloc)
	m.gauges["HeapIdle"] = float64(m.m.HeapIdle)
	m.gauges["HeapInuse"] = float64(m.m.HeapInuse)
	m.gauges["HeapObjects"] = float64(m.m.HeapObjects)
	m.gauges["HeapReleased"] = float64(m.m.HeapReleased)
	m.gauges["HeapSys"] = float64(m.m.HeapSys)
	m.gauges["LastGC"] = float64(m.m.LastGC)
	m.gauges["Lookups"] = float64(m.m.Lookups)
	m.gauges["MCacheInuse"] = float64(m.m.MCacheInuse)
	m.gauges["MCacheSys"] = float64(m.m.MCacheSys)
	m.gauges["MSpanInuse"] = float64(m.m.MSpanInuse)
	m.gauges["MSpanSys"] = float64(m.m.MSpanSys)
	m.gauges["Mallocs"] = float64(m.m.Mallocs)
	m.gauges["NextGC"] = float64(m.m.NextGC)
	m.gauges["NumForcedGC"] = float64(m.m.NumForcedGC)
	m.gauges["NumGC"] = float64(m.m.NumGC)
	m.gauges["OtherSys"] = float64(m.m.OtherSys)
	m.gauges["PauseTotalNs"] = float64(m.m.PauseTotalNs)
	m.gauges["StackInuse"] = float64(m.m.StackInuse)
	m.gauges["StackSys"] = float64(m.m.StackSys)
	m.gauges["Sys"] = float64(m.m.Sys)
	m.gauges["TotalAlloc"] = float64(m.m.TotalAlloc)

	// Custom metrics
	m.counters["PollCount"] += 1
	m.gauges["RandomValue"] = rand.Float64()

	// Polling
	if m.isStarted {
		time.Sleep(m.pollInterval)
		go m.update()
	}
}

func (m *RuntimeMonitor) Start() {
	m.isStarted = true
	go m.update()
}
func (m *RuntimeMonitor) Stop() {
	m.isStarted = false
}

func (m *RuntimeMonitor) GetCounters() map[string]int64 {
	return m.counters
}

func (m *RuntimeMonitor) GetGauges() map[string]float64 {
	return m.gauges
}
