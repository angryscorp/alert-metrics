package dbmetricstorage

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type PostgresMetricsStorage struct {
	pool *pgxpool.Pool
}

type Mock struct{}

func (m Mock) Ping() error {
	return errors.New("no ping for mock")
}

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

func (s PostgresMetricsStorage) Ping() error {
	return s.pool.Ping(context.TODO())
}
