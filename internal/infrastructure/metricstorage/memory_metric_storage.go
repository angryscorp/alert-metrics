package metricstorage

import (
	"errors"
	"github.com/angryscorp/alert-metrics/internal/domain"
	"sync"
)

var _ domain.MetricStorage = (*MemoryMetricStorage)(nil)

type MemoryMetricStorage struct {
	mu       sync.RWMutex
	gauges   map[string]float64
	counters map[string]int64
}

func NewMemoryMetricStorage() *MemoryMetricStorage {
	return &MemoryMetricStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}

func (m *MemoryMetricStorage) GetAllMetrics() []domain.Metric {
	m.mu.RLock()
	defer m.mu.RUnlock()

	res := make([]domain.Metric, 0)
	for key, value := range m.gauges {
		res = append(res, domain.Metric{
			ID:    key,
			MType: domain.MetricTypeGauge,
			Value: &value,
		})
	}
	for key, value := range m.counters {
		res = append(res, domain.Metric{
			ID:    key,
			MType: domain.MetricTypeCounter,
			Delta: &value,
		})
	}

	return res
}

func (m *MemoryMetricStorage) UpdateMetric(metrics domain.Metric) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	switch metrics.MType {
	case domain.MetricTypeCounter:
		m.counters[metrics.ID] += *metrics.Delta

	case domain.MetricTypeGauge:
		m.gauges[metrics.ID] = *metrics.Value

	default:
		return errors.New("unsupported metric type")
	}

	return nil
}

func (m *MemoryMetricStorage) GetMetric(metricType domain.MetricType, metricName string) (domain.Metric, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	res := domain.Metric{
		MType: metricType,
		ID:    metricName,
	}

	found := false
	switch metricType {
	case domain.MetricTypeCounter:
		v, ok := m.counters[metricName]
		if ok {
			res.Delta = &v
			found = true
		}

	case domain.MetricTypeGauge:
		v, ok := m.gauges[metricName]
		if ok {
			res.Value = &v
			found = true
		}
	}

	return res, found
}
