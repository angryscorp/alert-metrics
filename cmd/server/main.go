package main

import (
	"flag"
	"fmt"
	"github.com/angryscorp/alert-metrics/internal/config/server"
	"github.com/angryscorp/alert-metrics/internal/domain"
	"github.com/angryscorp/alert-metrics/internal/http/gzipper"
	"github.com/angryscorp/alert-metrics/internal/http/hash"
	"github.com/angryscorp/alert-metrics/internal/http/logger"
	"github.com/angryscorp/alert-metrics/internal/http/router"
	"github.com/angryscorp/alert-metrics/internal/infrastructure/dbmetricstorage"
	"github.com/angryscorp/alert-metrics/internal/infrastructure/metricstorage"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"os"
	"time"
)

func main() {
	config, err := server.NewConfig()
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err.Error())
		flag.Usage()
		os.Exit(1)
	}

	zeroLogger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	store, err := storeSelector(config, &zeroLogger)
	if err != nil {
		panic(err)
	}

	engine := gin.New()
	engine.
		Use(logger.New(zeroLogger)).
		Use(gin.Recovery()).
		Use(gzipper.UnzipMiddleware()).
		Use(hash.NewHashValidator(config.HashKey)).
		Use(gzip.Gzip(gzip.DefaultCompression))

	mr := router.New(engine, store)
	if err = mr.Run(config.Address); err != nil {
		panic(err)
	}
}

func storeSelector(config server.Config, logger *zerolog.Logger) (domain.MetricStorage, error) {
	if config.DatabaseDSN != "" {
		if err := dbmetricstorage.Migrate(config.DatabaseDSN); err != nil {
			return nil, fmt.Errorf("failed to migrate database: %w", err)
		}

		retryIntervals := []time.Duration{time.Second, time.Second * 3, time.Second * 5}
		var dbStore domain.MetricStorage
		dbStore, err := dbmetricstorage.New(config.DatabaseDSN, logger)
		if err != nil {
			return nil, err
		}

		return dbmetricstorage.NewRetryableDBStorage(dbStore, retryIntervals, logger), nil
	}

	if config.FileStoragePath != "" {
		return metricstorage.NewFileMetricStorage(
			metricstorage.NewMemoryMetricStorage(),
			*logger,
			time.Duration(config.StoreIntervalInSeconds)*time.Second,
			config.FileStoragePath,
			config.ShouldRestore,
		), nil
	}

	return metricstorage.NewMemoryMetricStorage(), nil
}
