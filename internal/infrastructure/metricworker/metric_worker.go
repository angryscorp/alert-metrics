package metricworker

import (
	"github.com/angryscorp/alert-metrics/internal/domain"
	"strconv"
	"time"
)

type MetricWorker struct {
	metricMonitor  domain.MetricMonitor
	metricReporter domain.MetricReporter
	reportInterval time.Duration
	isRunning      bool
}

func NewMetricWorker(metricMonitor domain.MetricMonitor, metricReporter domain.MetricReporter, reportInterval time.Duration) *MetricWorker {
	return &MetricWorker{
		metricMonitor:  metricMonitor,
		metricReporter: metricReporter,
		reportInterval: reportInterval,
	}
}

func (mw *MetricWorker) Start() {
	mw.isRunning = true
	go mw.sendCurrentMetrics()
}

func (mw *MetricWorker) sendCurrentMetrics() {
	metrics := mw.metricMonitor.GetMetrics()
	
	// Send Gauge metrics
	for key, value := range metrics.Gauges {
		formattedValue := strconv.FormatFloat(value, 'f', -1, 64)
		mw.metricReporter.Report(domain.MetricTypeGauge, key, formattedValue)
	}

	// Send Counter metrics
	for key, value := range metrics.Counters {
		formattedValue := strconv.FormatInt(value, 10)
		mw.metricReporter.Report(domain.MetricTypeCounter, key, formattedValue)
	}

	// Report interval
	if mw.isRunning {
		time.Sleep(mw.reportInterval)
		go mw.sendCurrentMetrics()
	}
}

func (mw *MetricWorker) Stop() {
	mw.isRunning = false
}
