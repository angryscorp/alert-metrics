package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/angryscorp/alert-metrics/internal/domain"
	"github.com/angryscorp/alert-metrics/internal/http/router"
)

type MetricsJSONHandler struct {
	storage domain.MetricStorage
}

func NewMetricsJSONHandler(storage domain.MetricStorage) MetricsJSONHandler {
	return MetricsJSONHandler{
		storage: storage,
	}
}

var _ router.MetricsJSONHandler = (*MetricsJSONHandler)(nil)

// FetchMetricsJSON handles a JSON POST request to fetch a specific metric by its type and ID from the storage.
func (handler MetricsJSONHandler) FetchMetricsJSON(c *gin.Context) {
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

// UpdateMetricsJSON handles an HTTP POST request to update a metric in the storage using the provided JSON payload.
func (handler MetricsJSONHandler) UpdateMetricsJSON(c *gin.Context) {
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

// BatchUpdateFetchMetrics handles a JSON request to batch update and fetch multiple metrics from storage.
func (handler MetricsJSONHandler) BatchUpdateFetchMetrics(c *gin.Context) {
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

func (handler MetricsJSONHandler) verifyContentTypeIsJSON(r *http.Request) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return errors.New("Content-Type must be application/json")
	}

	return nil
}

func (handler MetricsJSONHandler) updateMetrics(ctx context.Context, metrics domain.Metric) (domain.Metric, error) {
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
