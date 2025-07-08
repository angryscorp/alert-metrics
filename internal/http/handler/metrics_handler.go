package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/angryscorp/alert-metrics/internal/domain"
	"github.com/angryscorp/alert-metrics/internal/http/router"
)

type MetricsHandler struct {
	storage domain.MetricStorage
}

func New(storage domain.MetricStorage) MetricsHandler {
	return MetricsHandler{
		storage: storage,
	}
}

var _ router.MetricsHandler = (*MetricsHandler)(nil)

func (handler MetricsHandler) Ping(c *gin.Context) {
	if err := handler.storage.Ping(c.Request.Context()); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

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

func (handler MetricsHandler) UpdateMetrics(c *gin.Context) {
	if err := handler.update(c.Request.Context(), c.Param("metricType"), c.Param("metricName"), c.Param("metricValue")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.Status(http.StatusOK)
}

func (handler MetricsHandler) FetchMetricsJSON(c *gin.Context) {
	if err := handler.verifyContentTypeIsJSON(c.Request); err != nil {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": err.Error()})
		return
	}

	var metrics domain.Metric
	if err := c.ShouldBindJSON(&metrics); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, ok := handler.storage.GetMetric(c.Request.Context(), metrics.MType, metrics.ID)
	if ok {
		c.JSON(http.StatusOK, res)
		return
	}
	c.Header("Content-Type", "application/json")
	c.Status(http.StatusNotFound)
}

func (handler MetricsHandler) UpdateMetricsJSON(c *gin.Context) {
	if err := handler.verifyContentTypeIsJSON(c.Request); err != nil {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": err.Error()})
		return
	}

	var metric domain.Metric
	if err := c.ShouldBindJSON(&metric); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	metric, err := handler.updateMetrics(c.Request.Context(), metric)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, metric)
}

func (handler MetricsHandler) BatchUpdateFetchMetrics(c *gin.Context) {
	if err := handler.verifyContentTypeIsJSON(c.Request); err != nil {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": err.Error()})
		return
	}

	var metrics []domain.Metric
	if err := c.ShouldBindJSON(&metrics); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := handler.storage.UpdateMetrics(c.Request.Context(), metrics)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, metrics)
}

func (handler MetricsHandler) update(ctx context.Context, rawMetricType string, metricName string, metricValue string) error {
	metrics, err := domain.NewMetrics(rawMetricType, metricName, metricValue)
	if err != nil {
		return err
	}

	return handler.storage.UpdateMetric(ctx, *metrics)
}

func (handler MetricsHandler) updateMetrics(ctx context.Context, metrics domain.Metric) (domain.Metric, error) {
	err := handler.storage.UpdateMetric(ctx, metrics)
	if err != nil {
		return domain.Metric{}, err
	}

	res, ok := handler.storage.GetMetric(ctx, metrics.MType, metrics.ID)
	if !ok {
		return domain.Metric{}, errors.New("failed to get updated metrics")
	}
	return res, nil
}

func (handler MetricsHandler) verifyContentTypeIsJSON(r *http.Request) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return errors.New("Content-Type must be application/json")
	}

	return nil
}
