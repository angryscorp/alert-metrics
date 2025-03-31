package metricrouter

import (
	"github.com/angryscorp/alert-metrics/internal/domain"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

type MetricRouter struct {
	mux    *http.ServeMux
	store  domain.MetricStorage
	logger *slog.Logger
}

func NewMetricRouter(mux *http.ServeMux, store domain.MetricStorage) MetricRouter {
	router := MetricRouter{}
	router.mux = mux
	router.store = store
	router.logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	router.registerHealthCheck()
	router.registerMetricUpdate()
	return router
}

func (router MetricRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.mux.ServeHTTP(w, r)
}

func (router MetricRouter) registerHealthCheck() {
	router.mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		router.logger.Info("health check request", "method", r.Method, "url", r.URL.String())
		_, err := w.Write([]byte("OK"))
		if err != nil {
			router.logger.Error("health check response writing failed", "method", r.Method, "url", r.URL.String())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		router.logger.Info("health check response is ok", "method", r.Method, "url", r.URL.String())
		w.WriteHeader(http.StatusOK)
	})
}

// The function is expected only POST request with the following format: /update/metric_type/metric_name/metric_value
// Metric type is enum with the following options: counter, gauge.
func (router MetricRouter) registerMetricUpdate() {
	allowedPath := "update"
	allowedMethod := http.MethodPost
	allowedContentType := "text/plain"

	router.mux.HandleFunc("/update/", func(w http.ResponseWriter, r *http.Request) {
		router.logRequest("update metric request", r)

		// Check for allowed method
		if r.Method != allowedMethod {
			router.logError("method is not allowed", r, nil)
			http.Error(w, "invalid request method", http.StatusBadRequest)
			return
		}

		// Check for allowed content type
		if r.Header.Get("Content-Type") != allowedContentType {
			router.logError("content-type is not valid", r, nil)
			http.Error(w, "invalid content-type", http.StatusBadRequest)
			return
		}

		// Check URL structure
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) != 4 || parts[0] != allowedPath {
			router.logError("invalid path or url format", r, nil)
			http.Error(w, "invalid path or URL format", http.StatusNotFound)
			return
		}

		// Check for Metric type
		metricType, err := domain.NewMetricType(parts[1])
		if err != nil {
			router.logError("invalid metric type", r, err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Update metrics in the storage
		metricName, metricValue := parts[2], parts[3]
		if err := router.store.Update(metricType, metricName, metricValue); err != nil {
			router.logError("failed to update metric", r, err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// All good, metrics are updated
		router.logger.Info("metrics have been updated", "metric type", metricType, "metric name", metricName, "metric value", metricValue)

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
	})
}

func (router MetricRouter) logRequest(msg string, r *http.Request) {
	router.logger.Info(msg, "method", r.Method, "url", r.URL.String(), "content-type", r.Header.Get("Content-Type"))
}

func (router MetricRouter) logError(msg string, r *http.Request, err error) {
	errMsg := safeErrorMsg(err, msg)
	router.logger.Error(msg, "method", r.Method, "url", r.URL.String(), "content-type", r.Header.Get("Content-Type"), "error", errMsg)
}

func safeErrorMsg(err error, fallback string) string {
	if err != nil {
		return err.Error()
	}
	return fallback
}
