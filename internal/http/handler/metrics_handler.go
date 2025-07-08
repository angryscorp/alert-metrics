package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/angryscorp/alert-metrics/internal/domain"
	"github.com/angryscorp/alert-metrics/internal/http/router"
)

type MetricsHandler struct {
	storage domain.MetricStorage
}

func NewMetricsHandler(storage domain.MetricStorage) MetricsHandler {
	return MetricsHandler{
		storage: storage,
	}
}

var _ router.MetricsHandler = (*MetricsHandler)(nil)

// GetMetric retrieves the specified metric from storage and returns its value in the response. Responds with appropriate status codes.
func (handler MetricsHandler) GetMetric(c *gin.Context) {
	metricType, err := domain.NewMetricType(c.Param("metricType"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	metrics, ok := handler.storage.GetMetric(c.Request.Context(), metricType, c.Param("metricName"))
	if !ok {
		c.Status(http.StatusNotFound)
		return
	}

	c.String(http.StatusOK, metrics.StringValue())
}

// GetAllMetrics retrieves all stored metrics and returns them as an HTML response. Shows a "No data" message if no metrics exist.
func (handler MetricsHandler) GetAllMetrics(c *gin.Context) {
	allMetrics := handler.storage.GetAllMetrics(c.Request.Context())
	if len(allMetrics) == 0 {
		c.Data(http.StatusOK, "text/html", []byte("<h3>No data</h3>"))
		return
	}

	htmlContent := "<h3>Current metrics</h3><ul>"
	for _, v := range domain.NewMetricRepresentatives(allMetrics).SortByName() {
		htmlContent += fmt.Sprintf("<li>%s</li>", v)
	}
	htmlContent += "</ul>"
	c.Data(http.StatusOK, "text/html", []byte(htmlContent))
}

// UpdateMetrics processes an update request for a specific metric and responds with an appropriate HTTP status code.
func (handler MetricsHandler) UpdateMetrics(c *gin.Context) {
	if err := handler.update(c.Request.Context(), c.Param("metricType"), c.Param("metricName"), c.Param("metricValue")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.Status(http.StatusOK)
}

func (handler MetricsHandler) update(ctx context.Context, rawMetricType string, metricName string, metricValue string) error {
	metrics, err := domain.NewMetrics(rawMetricType, metricName, metricValue)
	if err != nil {
		return err
	}

	return handler.storage.UpdateMetric(ctx, *metrics)
}
