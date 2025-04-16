package metricstorage

import (
	"errors"
	"github.com/angryscorp/alert-metrics/internal/domain"
	"strconv"
	"sync"
)

var _ domain.MetricStorage = (*MemStorage)(nil)

type MemStorage struct {
	mu       sync.RWMutex
	gauges   map[string]float64
	counters map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}

func (m *MemStorage) GetAllMetrics() map[string]string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	res := make(map[string]string)
	for key, value := range m.gauges {
		res[string(domain.MetricTypeGauge)+"."+key] = strconv.FormatFloat(value, 'f', -1, 64)
	}
	for key, value := range m.counters {
		res[string(domain.MetricTypeCounter)+"."+key] = strconv.FormatInt(value, 10)
	}
	return res
}

func (m *MemStorage) UpdateMetrics(metrics domain.Metrics) error {
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

func (m *MemStorage) GetMetrics(metricType domain.MetricType, metricName string) (domain.Metrics, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	res := domain.Metrics{
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
