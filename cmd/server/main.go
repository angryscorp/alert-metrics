package main

import (
	"flag"
	"fmt"
	"github.com/angryscorp/alert-metrics/internal/domain"
	"github.com/angryscorp/alert-metrics/internal/infrastructure/httplogger"
	"github.com/angryscorp/alert-metrics/internal/infrastructure/metricdumper"
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

	var initData *[]domain.Metric
	if config.ShouldRestore {
		initData = nil // read from the file to initData (config.fileStoragePath)
	}

	var store domain.MetricStorage
	store, err = metricstorage.New(initData)
	if err != nil {
		panic("initializing metric storage failed: " + err.Error())
	}

	if config.FileStoragePath != "" {
		writer := os.Stdout // init a new writer to the file (config.fileStoragePath)
		store = metricdumper.New(store, time.Duration(config.StoreIntervalInSeconds), writer, logger)
	}

	var mr = metricrouter.NewMetricRouter(
		router,
		store,
	)

	err = mr.Run(config.Address)
	if err != nil {
		panic(err)
	}
}
