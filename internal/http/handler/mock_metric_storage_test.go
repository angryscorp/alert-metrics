package handler

import (
	"context"
	"github.com/angryscorp/alert-metrics/internal/domain"
	"github.com/stretchr/testify/mock"
)

// MockMetricStorage is a mock implementation of domain.MetricStorage for testing
type MockMetricStorage struct {
	mock.Mock
}

func (m *MockMetricStorage) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockMetricStorage) GetAllMetrics(ctx context.Context) []domain.Metric {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Metric)
}

func (m *MockMetricStorage) UpdateMetric(ctx context.Context, metric domain.Metric) error {
	args := m.Called(ctx, metric)
	return args.Error(0)
}

func (m *MockMetricStorage) UpdateMetrics(ctx context.Context, metrics []domain.Metric) error {
	args := m.Called(ctx, metrics)
	return args.Error(0)
}

func (m *MockMetricStorage) GetMetric(ctx context.Context, metricType domain.MetricType, metricName string) (domain.Metric, bool) {
	args := m.Called(ctx, metricType, metricName)
	return args.Get(0).(domain.Metric), args.Bool(1)
}
