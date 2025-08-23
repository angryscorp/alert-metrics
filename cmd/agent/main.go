package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/angryscorp/alert-metrics/internal/buildinfo"
	"github.com/angryscorp/alert-metrics/internal/crypto"
	cryptohttp "github.com/angryscorp/alert-metrics/internal/http/crypto"
	"github.com/angryscorp/alert-metrics/internal/http/realip"
	"github.com/angryscorp/alert-metrics/internal/infrastructure/shutdown"

	"github.com/rs/zerolog"

	"github.com/angryscorp/alert-metrics/internal/config/agent"
	"github.com/angryscorp/alert-metrics/internal/http/gzipper"
	"github.com/angryscorp/alert-metrics/internal/http/hash"
	"github.com/angryscorp/alert-metrics/internal/http/retry"
	"github.com/angryscorp/alert-metrics/internal/infrastructure/metricmonitor"
	"github.com/angryscorp/alert-metrics/internal/infrastructure/metricreporter"
	"github.com/angryscorp/alert-metrics/internal/infrastructure/metricworker"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	fmt.Printf("%s", buildinfo.New(buildVersion, buildDate, buildCommit))

	flags, err := agent.NewConfig()
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err.Error())
		flag.Usage()
		log.Fatal(err.Error())
	}

	runtimeMonitor := metricmonitor.NewRuntimeMonitor(time.Duration(flags.PollIntervalInSeconds) * time.Second)
	runtimeMonitor.Start()

	metricReporter := metricreporter.NewHTTPMetricReporter(
		"http://"+flags.Address,
		&http.Client{
			Transport: buildTransport(
				flags.PathToCryptoKey,
				flags.HashKey,
				[]time.Duration{time.Second, time.Second * 3, time.Second * 5},
				zerolog.New(os.Stdout).With().Timestamp().Logger(),
			),
		},
	)

	worker := metricworker.NewMetricWorker(
		runtimeMonitor,
		metricReporter,
		time.Duration(flags.ReportIntervalInSeconds)*time.Second,
		flags.RateLimit,
	)

	worker.RunWithGracefulShutdown(shutdown.NewGracefulShutdownNotifier())

	select {}
}

func buildTransport(cryptoKeyPath, hashKey string, retryIntervals []time.Duration, logger zerolog.Logger) http.RoundTripper {
	// Base transport
	transport := http.DefaultTransport

	// Real IP transport
	transport = realip.New(transport)

	// Gzip transport
	transport = gzipper.NewGzipTransport(transport)

	// Crypto transport
	if cryptoKeyPath != "" {
		encryptor, err := crypto.NewPublicKeyEncrypter(cryptoKeyPath)
		if err != nil {
			log.Fatal(err.Error())
		}
		transport = cryptohttp.EncryptorMiddleware(encryptor)(transport)
	}

	// Hash transport
	transport = hash.NewHashTransport(transport, hashKey)

	// Retry transport
	transport = retry.New(transport, retryIntervals, logger)

	return transport
}
