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

	router := gin.New()
	router.
		Use(httplogger.New(logger)).
		Use(gin.Recovery()).
		Use(gzip.Gzip(gzip.DefaultCompression))

	var store domain.MetricStorage
	store = metricstorage.NewMemoryMetricStorage()

	if config.FileStoragePath != "" {
		store = metricstorage.NewFileMetricStorage(
			store,
			logger,
			time.Duration(config.StoreIntervalInSeconds)*time.Second,
			config.FileStoragePath,
			config.ShouldRestore,
		)
	}

	var db interface{ Ping() error }
	if config.DatabaseDSN != "" {
		db, err = dbmetricstorage.New(config.DatabaseDSN)
		if err != nil {
			panic(err)
		}
	} else {
		db = dbmetricstorage.Mock{}
	}

	var mr = metricrouter.NewMetricRouter(
		router,
		store,
		db,
	)

	err = mr.Run(config.Address)
	if err != nil {
		panic(err)
	}
}
