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

func (m *MemStorage) Get(metricType domain.MetricType, key string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	switch metricType {
	case domain.MetricTypeCounter:
		res, ok := m.counters[key]
		if ok {
			return strconv.FormatInt(res, 10), true
		}

	case domain.MetricTypeGauge:
		res, ok := m.gauges[key]
		if ok {
			return strconv.FormatFloat(res, 'f', -1, 64), true
		}
	}

	return "", false
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

func (m *MemStorage) Update(metricType domain.MetricType, key string, value string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if key == "" {
		return errors.New("metric name is empty")
	}

	switch metricType {
	case domain.MetricTypeCounter:
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return errors.New("invalid counter value")
		}
		m.counters[key] += v

	case domain.MetricTypeGauge:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return errors.New("invalid gauge value")
		}
		m.gauges[key] = v

	default:
		return errors.New("unsupported metric type")
	}

	return nil
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

func (m *MemStorage) GetMetrics(metricType domain.MetricType, metricName string) domain.Metrics {
	m.mu.RLock()
	defer m.mu.RUnlock()

	res := domain.Metrics{
		MType: metricType,
		ID:    metricName,
	}

	switch metricType {
	case domain.MetricTypeCounter:
		v := m.counters[metricName]
		res.Delta = &v

	case domain.MetricTypeGauge:
		v := m.gauges[metricName]
		res.Value = &v
	}

	return res
}
