package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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
	showBuildInfo()

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
			Transport: retry.New(
				hash.NewHashTransport(
					gzipper.NewGzipTransport(http.DefaultTransport),
					flags.HashKey,
				),
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
	worker.Start()

	select {}
}

func showBuildInfo() {
	if buildVersion == "" {
		buildVersion = "N/A"
	}

	if buildDate == "" {
		buildDate = "N/A"
	}

	if buildCommit == "" {
		buildCommit = "N/A"
	}

	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)
}
