package main

import (
	"flag"
	"fmt"
	"github.com/angryscorp/alert-metrics/internal/domain"
	"github.com/angryscorp/alert-metrics/internal/infrastructure/dbmetricstorage"
	"github.com/angryscorp/alert-metrics/internal/infrastructure/httplogger"
	"github.com/angryscorp/alert-metrics/internal/infrastructure/metricrouter"
	"github.com/angryscorp/alert-metrics/internal/infrastructure/metricstorage"
	"github.com/angryscorp/alert-metrics/internal/infrastructure/serverconfig"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"os"
	"time"
)

func main() {
	config, err := serverconfig.New()
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err.Error())
		flag.Usage()
		os.Exit(1)
	}

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	store, err := storeSelector(config, &logger)
	if err != nil {
		panic(err)
	}

	router := gin.New()
	router.
		Use(httplogger.New(logger)).
		Use(gin.Recovery()).
		Use(gzip.Gzip(gzip.DefaultCompression))

	mr := metricrouter.New(router, store)
	if err = mr.Run(config.Address); err != nil {
		panic(err)
	}
}

func storeSelector(config serverconfig.ServerConfig, logger *zerolog.Logger) (domain.MetricStorage, error) {
	if config.DatabaseDSN != "" {
		return dbmetricstorage.New(config.DatabaseDSN, logger)
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
