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

func (m *MemStorage) Get(metricType domain.MetricType, key string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	switch metricType {
	case domain.MetricTypeCounter:
		res, ok := m.counters[key]
		if ok {
			return strconv.FormatInt(res, 10), nil
		}
		return "", errors.New("not found")

	case domain.MetricTypeGauge:
		res, ok := m.gauges[key]
		if ok {
			return strconv.FormatFloat(res, 'f', -1, 64), nil
		}
		return "", errors.New("not found")

	default:
		return "", errors.New("unsupported metric type")
	}
}

func (m *MemStorage) Update(metricType domain.MetricType, key string, value string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

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
