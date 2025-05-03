package dbmetricstorage

import (
	"github.com/angryscorp/alert-metrics/internal/domain"
	"github.com/rs/zerolog"
	"time"
)

type RetryablePostgresStorage struct {
	storage        domain.MetricStorage
	retryIntervals []time.Duration
	logger         *zerolog.Logger
}

var _ domain.MetricStorage = (*RetryablePostgresStorage)(nil)

func NewRetryableDBStorage(
	storage domain.MetricStorage,
	retryAttempts []time.Duration,
	logger *zerolog.Logger,
) *RetryablePostgresStorage {
	return &RetryablePostgresStorage{
		storage:        storage,
		retryIntervals: append([]time.Duration{0}, retryAttempts...),
		logger:         logger,
	}
}

func (s *RetryablePostgresStorage) UpdateMetric(metric domain.Metric) error {
	return s.withRetry(func() error {
		return s.storage.UpdateMetric(metric)
	})
}

func (s *RetryablePostgresStorage) GetAllMetrics() []domain.Metric {
	return s.storage.GetAllMetrics()
}

func (s *RetryablePostgresStorage) UpdateMetrics(metrics []domain.Metric) error {
	return s.withRetry(func() error {
		return s.storage.UpdateMetrics(metrics)
	})
}

func (s *RetryablePostgresStorage) GetMetric(metricType domain.MetricType, metricName string) (domain.Metric, bool) {
	return s.storage.GetMetric(metricType, metricName)
}

func (s *RetryablePostgresStorage) Ping() error {
	return s.withRetry(func() error {
		return s.storage.Ping()
	})
}

func (s *RetryablePostgresStorage) withRetry(target func() error) error {
	var err error
	for n, interval := range s.retryIntervals {
		if n > 0 {
			s.logger.Warn().Msgf("retrying connect %d in %f", n, interval.Seconds())
			time.Sleep(interval)
		}

		err = target()
		if err == nil {
			return nil
		}
		s.logger.Error().Err(err).Msgf("failed to connect %s", err)
	}

	return err
}
