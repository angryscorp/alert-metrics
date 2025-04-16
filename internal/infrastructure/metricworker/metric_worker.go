package metricworker

import (
	"github.com/angryscorp/alert-metrics/internal/domain"
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
		m := domain.Metrics{
			ID:    key,
			MType: domain.MetricTypeGauge,
			Value: &value,
		}
		mw.metricReporter.ReportMetrics(m)
	}

	// Send Counter metrics
	for key, value := range metrics.Counters {
		m := domain.Metrics{
			ID:    key,
			MType: domain.MetricTypeCounter,
			Delta: &value,
		}
		mw.metricReporter.ReportMetrics(m)
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
