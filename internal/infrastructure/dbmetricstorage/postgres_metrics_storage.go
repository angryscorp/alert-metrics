package dbmetricstorage

import (
	"context"
	"errors"
	"fmt"
	"github.com/angryscorp/alert-metrics/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"time"
)

type PostgresMetricsStorage struct {
	pool   *pgxpool.Pool
	logger *zerolog.Logger
}

var _ domain.MetricStorage = (*PostgresMetricsStorage)(nil)

func New(dsn string, logger *zerolog.Logger) (*PostgresMetricsStorage, error) {
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

	store := &PostgresMetricsStorage{pool: pool, logger: logger}
	if err := store.prepareDataTable(); err != nil {
		return nil, fmt.Errorf("failed to prepare data table: %w", err)
	}

	return store, nil
}

func (s PostgresMetricsStorage) GetAllMetrics() []domain.Metric {
	metrics := make([]domain.Metric, 0)

	rows, err := s.pool.Query(context.TODO(), selectAllMetrics)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to query metrics")
		return metrics
	}

	defer rows.Close()

	for rows.Next() {
		var metric domain.Metric
		err := rows.Scan(&metric.ID, &metric.MType, &metric.Delta, &metric.Value)
		if err != nil {
			panic(err)
		}
		metrics = append(metrics, metric)
	}

	if err := rows.Err(); err != nil {
		s.logger.Error().Err(err).Msg("failed to scan metrics")
	}

	return metrics
}

func (s PostgresMetricsStorage) UpdateMetric(metric domain.Metric) error {
	_, err := s.pool.Exec(context.TODO(), upsertMetric,
		metric.ID, metric.MType, metric.Delta, metric.Value,
	)

	if err != nil {
		return fmt.Errorf("failed to update metric: %w", err)
	}

	return nil
}

func (s PostgresMetricsStorage) UpdateMetrics(metrics []domain.Metric) error {
	ctx := context.TODO()
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	for _, metric := range metrics {
		_, err = tx.Exec(ctx, upsertMetric,
			metric.ID, metric.MType, metric.Delta, metric.Value,
		)
		if err != nil {
			return fmt.Errorf("failed to update metric: %w", err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s PostgresMetricsStorage) GetMetric(metricType domain.MetricType, metricName string) (domain.Metric, bool) {
	row := s.pool.QueryRow(context.TODO(), selectMetric,
		metricName, metricType,
	)

	metric := domain.Metric{ID: metricName, MType: metricType}
	err := row.Scan(&metric.Delta, &metric.Value)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Metric{}, false
		}
	}

	return metric, true
}

func (s PostgresMetricsStorage) Ping() error {
	return s.pool.Ping(context.TODO())
}

func (s PostgresMetricsStorage) prepareDataTable() error {
	_, err := s.pool.Exec(context.TODO(), createTableMetrics)

	return err
}
