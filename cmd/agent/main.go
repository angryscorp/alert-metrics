package main

import (
	"github.com/angryscorp/alert-metrics/internal/infrastructure/metricmonitor"
	"github.com/angryscorp/alert-metrics/internal/infrastructure/metricreporter"
	"github.com/angryscorp/alert-metrics/internal/infrastructure/metricworker"
	"net/http"
	"time"
)

func main() {
	rm := metricmonitor.NewRuntimeMonitor(2 * time.Second)
	rm.Start()

	mr := metricreporter.NewHTTPMetricReporter("http://localhost:8080", &http.Client{})

	worker := metricworker.NewMetricWorker(rm, mr, 10*time.Second)
	worker.Start()

	select {}
}
