package dbmetricstorage

import (
	"context"
	"fmt"
	"github.com/angryscorp/alert-metrics/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type PostgresMetricsStorage struct {
	pool *pgxpool.Pool
}

var _ domain.MetricStorage = (*PostgresMetricsStorage)(nil)

func New(dsn string) (*PostgresMetricsStorage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse dsn: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresMetricsStorage{pool: pool}, nil
}

func (s PostgresMetricsStorage) GetAllMetrics() []domain.Metric {
	//TODO implement me
	//panic("implement me")
	return []domain.Metric{}
}

func (s PostgresMetricsStorage) UpdateMetric(metric domain.Metric) error {
	//TODO implement me
	//panic("implement me")
	return nil
}

func (s PostgresMetricsStorage) GetMetric(metricType domain.MetricType, metricName string) (domain.Metric, bool) {
	//TODO implement me
	//panic("implement me")
	return domain.Metric{}, false
}

func (s PostgresMetricsStorage) Ping() error {
	return s.pool.Ping(context.TODO())
}
