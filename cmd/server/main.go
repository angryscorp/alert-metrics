package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/angryscorp/alert-metrics/internal/http/subnet"

	"github.com/angryscorp/alert-metrics/internal/infrastructure/shutdown"

	grpcserver "github.com/angryscorp/alert-metrics/internal/grpc/server"

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

	shutdownCh := shutdown.NewGracefulShutdownNotifier()
	serverCount := 1 // HTTP always running
	if config.UseGRPC {
		serverCount = 2 // + gRPC server
	}

	var wg sync.WaitGroup
	errChan := make(chan error, serverCount)

	// HTTP server always running
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := runHTTPServer(config, store, zeroLogger, shutdownCh); err != nil {
			errChan <- fmt.Errorf("HTTP server error: %w", err)
		}
	}()

	// gRPC server is optional
	if config.UseGRPC {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := runGRPCServer(config, store, zeroLogger, shutdownCh); err != nil {
				errChan <- fmt.Errorf("gRPC server error: %w", err)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			log.Printf("Server error: %v", err)
		}
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

func runHTTPServer(config server.Config, store domain.MetricStorage, zeroLogger zerolog.Logger, shutdownCh <-chan struct{}) error {
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
			return fmt.Errorf("failed to create decrypter: %w", err)
		}
		engine.Use(cryptohttp.DecrypterMiddleware(decrypter))
	}

	mr := router.New(engine, &zeroLogger)
	mr.RegisterPingHandler(handler.NewPingHandler(store))
	mr.RegisterMetricsHandler(handler.NewMetricsHandler(store))
	mr.RegisterMetricsJSONHandler(handler.NewMetricsJSONHandler(store))

	zeroLogger.Info().Str("address", config.Address).Msg("starting HTTP server")
	return mr.Run(config.Address, shutdownCh)
}

func runGRPCServer(config server.Config, store domain.MetricStorage, zeroLogger zerolog.Logger, shutdownCh <-chan struct{}) error {
	grpcSrv := grpcserver.NewGRPCServer(store, zeroLogger)

	zeroLogger.Info().Str("address", config.GRPCAddress).Msg("starting gRPC server")
	return grpcSrv.Run(config.GRPCAddress, shutdownCh)
}
