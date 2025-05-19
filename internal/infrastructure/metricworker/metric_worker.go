package metricworker

import (
	"github.com/angryscorp/alert-metrics/internal/domain"
	"time"
)

const batchSize = 10

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
	go mw.sendBatch()
}

func (mw *MetricWorker) sendCurrentMetrics() {
	metrics := mw.metricMonitor.GetMetrics()

	// Send Gauge metrics
	for key, value := range metrics.Gauges {
		m := domain.Metric{
			ID:    key,
			MType: domain.MetricTypeGauge,
			Value: &value,
		}
		mw.metricReporter.ReportMetric(m)
	}

	// Send Counter metrics
	for key, value := range metrics.Counters {
		m := domain.Metric{
			ID:    key,
			MType: domain.MetricTypeCounter,
			Delta: &value,
		}
		mw.metricReporter.ReportMetric(m)
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

func (mw *MetricWorker) sendBatch() {
	buf := make([]domain.Metric, 0)
	rawMetrics := mw.metricMonitor.GetMetrics()

	// Send Gauge metrics
	for key, value := range rawMetrics.Gauges {
		metric := domain.Metric{
			ID:    key,
			MType: domain.MetricTypeGauge,
			Value: &value,
		}
		buf = append(buf, metric)
		if len(buf) >= batchSize {
			mw.metricReporter.ReportBatch(buf)
			buf = make([]domain.Metric, 0)
		}
	}

	// Send Counter metrics
	for key, value := range rawMetrics.Counters {
		metric := domain.Metric{
			ID:    key,
			MType: domain.MetricTypeCounter,
			Delta: &value,
		}
		buf = append(buf, metric)
		if len(buf) >= batchSize {
			mw.metricReporter.ReportBatch(buf)
			buf = make([]domain.Metric, 0)
		}
	}

	if len(buf) > 0 {
		mw.metricReporter.ReportBatch(buf)
	}

	// Report interval
	if mw.isRunning {
		time.Sleep(mw.reportInterval)
		go mw.sendBatch()
	}
}
