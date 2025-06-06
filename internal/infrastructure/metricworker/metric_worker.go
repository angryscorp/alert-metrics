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
	rateLimiter    int
	isRunning      bool
	requestChan    chan []domain.Metric
}

func NewMetricWorker(
	metricMonitor domain.MetricMonitor,
	metricReporter domain.MetricReporter,
	reportInterval time.Duration,
	rateLimiter int,
) *MetricWorker {
	return &MetricWorker{
		metricMonitor:  metricMonitor,
		metricReporter: metricReporter,
		reportInterval: reportInterval,
		rateLimiter:    rateLimiter,
		requestChan:    make(chan []domain.Metric),
	}
}

func (mw *MetricWorker) Start() {
	mw.isRunning = true
	go mw.sendBatch()
	go mw.startWorkerPool()
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
			mw.requestChan <- buf
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
			mw.requestChan <- buf
			buf = make([]domain.Metric, 0)
		}
	}

	if len(buf) > 0 {
		mw.requestChan <- buf
	}

	// Report interval
	if mw.isRunning {
		time.Sleep(mw.reportInterval)
		go mw.sendBatch()
	}
}

func (mw *MetricWorker) startWorkerPool() {
	for i := 0; i < mw.rateLimiter; i++ {
		go func(ch chan []domain.Metric) {
			for req := range ch {
				mw.metricReporter.ReportBatch(req)
			}
		}(mw.requestChan)
	}
}
