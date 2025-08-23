package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/angryscorp/alert-metrics/internal/http/subnet"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/angryscorp/alert-metrics/internal/infrastructure/shutdown"

	"github.com/angryscorp/alert-metrics/internal/buildinfo"
	"github.com/angryscorp/alert-metrics/internal/config/server"
	"github.com/angryscorp/alert-metrics/internal/crypto"
	"github.com/angryscorp/alert-metrics/internal/domain"
	cryptohttp "github.com/angryscorp/alert-metrics/internal/http/crypto"
	"github.com/angryscorp/alert-metrics/internal/http/gzipper"
	"github.com/angryscorp/alert-metrics/internal/http/handler"
	"github.com/angryscorp/alert-metrics/internal/http/hash"
	"github.com/angryscorp/alert-metrics/internal/http/logger"
	"github.com/angryscorp/alert-metrics/internal/http/router"
	"github.com/angryscorp/alert-metrics/internal/infrastructure/dbmetricstorage"
	"github.com/angryscorp/alert-metrics/internal/infrastructure/metricstorage"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	fmt.Printf("%s", buildinfo.New(buildVersion, buildDate, buildCommit))

	config, err := server.NewConfig()
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err.Error())
		flag.Usage()
		log.Fatal(err.Error())
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
		Use(subnet.NewTrustedSubnetMiddleware(config.TrustedSubnet)).
		Use(gzip.Gzip(gzip.DefaultCompression))

	if config.PathToCryptoKey != "" {
		decrypter, err := crypto.NewPrivateKeyDecrypter(config.PathToCryptoKey)
		if err != nil {
			log.Fatal(err.Error())
		}
		engine.Use(cryptohttp.DecrypterMiddleware(decrypter))
	}

	mr := router.New(engine, &zeroLogger)
	mr.RegisterPingHandler(handler.NewPingHandler(store))
	mr.RegisterMetricsHandler(handler.NewMetricsHandler(store))
	mr.RegisterMetricsJSONHandler(handler.NewMetricsJSONHandler(store))

	shutdownCh := shutdown.NewGracefulShutdownNotifier()
	if err = mr.Run(config.Address, shutdownCh); err != nil {
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
