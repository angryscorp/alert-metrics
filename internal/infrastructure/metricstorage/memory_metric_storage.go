package metricstorage

import (
	"errors"
	"github.com/angryscorp/alert-metrics/internal/domain"
	"strconv"
	"sync"
)

var _ domain.MetricStorage = (*MemoryMetricStorage)(nil)

type MemoryMetricStorage struct {
	mu       sync.RWMutex
	gauges   map[string]float64
	counters map[string]int64
}

func New(initialData *[]domain.Metric) (*MemoryMetricStorage, error) {
	store := &MemoryMetricStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}

	if initialData == nil {
		return store, nil
	}

	for _, v := range *initialData {
		if err := store.UpdateMetrics(v); err != nil {
			return nil, err
		}
	}
	return store, nil
}

func (m *MemoryMetricStorage) GetAllMetrics() domain.MetricRepresentatives {
	m.mu.RLock()
	defer m.mu.RUnlock()

	res := make([]domain.MetricRepresentative, 0)
	for key, value := range m.gauges {
		res = append(res, domain.MetricRepresentative{
			Type:  domain.MetricTypeGauge,
			Name:  key,
			Value: strconv.FormatFloat(value, 'f', -1, 64),
		})
	}
	for key, value := range m.counters {
		res = append(res, domain.MetricRepresentative{
			Type:  domain.MetricTypeCounter,
			Name:  key,
			Value: strconv.FormatInt(value, 10),
		})
	}

	return res
}

func (m *MemoryMetricStorage) UpdateMetrics(metrics domain.Metric) error {
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

func (m *MemoryMetricStorage) GetMetrics(metricType domain.MetricType, metricName string) (domain.Metric, bool) {
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
