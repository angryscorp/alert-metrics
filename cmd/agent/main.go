package main

import (
	"flag"
	"fmt"
	"github.com/angryscorp/alert-metrics/internal/infrastructure/agentconfig"
	"github.com/angryscorp/alert-metrics/internal/infrastructure/metricmonitor"
	"github.com/angryscorp/alert-metrics/internal/infrastructure/metricreporter"
	"github.com/angryscorp/alert-metrics/internal/infrastructure/metricworker"
	"github.com/angryscorp/alert-metrics/internal/infrastructure/zipper"
	"net/http"
	"os"
	"time"
)

func main() {
	flags, err := agentconfig.NewAgentConfig()
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
			Transport: zipper.NewGzipTransport(http.DefaultTransport),
		},
	)

	worker := metricworker.NewMetricWorker(rm, mr, time.Duration(flags.ReportIntervalInSeconds)*time.Second)
	worker.Start()

	select {}
}
