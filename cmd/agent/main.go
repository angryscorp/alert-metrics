package main

import (
	"flag"
	"fmt"
	"github.com/angryscorp/alert-metrics/internal/infrastructure/agentconfig"
	"github.com/angryscorp/alert-metrics/internal/infrastructure/gzipper"
	"github.com/angryscorp/alert-metrics/internal/infrastructure/hash"
	"github.com/angryscorp/alert-metrics/internal/infrastructure/metricmonitor"
	"github.com/angryscorp/alert-metrics/internal/infrastructure/metricreporter"
	"github.com/angryscorp/alert-metrics/internal/infrastructure/metricworker"
	"github.com/angryscorp/alert-metrics/internal/infrastructure/retrytransport"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"time"
)

func main() {
	flags, err := agentconfig.New()
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err.Error())
		flag.Usage()
		os.Exit(1)
	}

	rm := metricmonitor.NewRuntimeMonitor(time.Duration(flags.PollIntervalInSeconds) * time.Second)
	rm.Start()

	mr := metricreporter.NewHTTPMetricReporter(
		"http://"+flags.Address,
		&http.Client{
			Transport: retrytransport.New(
				hash.NewHashTransport(
					gzipper.NewGzipTransport(http.DefaultTransport),
					flags.HashKey,
				),
				[]time.Duration{time.Second, time.Second * 3, time.Second * 5},
				zerolog.New(os.Stdout).With().Timestamp().Logger(),
			),
		},
	)

	worker := metricworker.NewMetricWorker(rm, mr, time.Duration(flags.ReportIntervalInSeconds)*time.Second)
	worker.Start()

	select {}
}
