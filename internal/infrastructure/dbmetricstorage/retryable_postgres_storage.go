package dbmetricstorage

import (
	"context"
	"time"

	"github.com/rs/zerolog"

	"github.com/angryscorp/alert-metrics/internal/domain"
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

func (s *RetryablePostgresStorage) UpdateMetric(ctx context.Context, metric domain.Metric) error {
	return s.withRetry(func() error {
		return s.storage.UpdateMetric(ctx, metric)
	})
}

func (s *RetryablePostgresStorage) GetAllMetrics(ctx context.Context) []domain.Metric {
	return s.storage.GetAllMetrics(ctx)
}

func (s *RetryablePostgresStorage) UpdateMetrics(ctx context.Context, metrics []domain.Metric) error {
	return s.withRetry(func() error {
		return s.storage.UpdateMetrics(ctx, metrics)
	})
}

func (s *RetryablePostgresStorage) GetMetric(ctx context.Context, metricType domain.MetricType, metricName string) (domain.Metric, bool) {
	return s.storage.GetMetric(ctx, metricType, metricName)
}

func (s *RetryablePostgresStorage) Ping(ctx context.Context) error {
	return s.withRetry(func() error {
		return s.storage.Ping(ctx)
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
