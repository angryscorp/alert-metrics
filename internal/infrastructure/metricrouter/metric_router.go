package metricrouter

import (
	"errors"
	"fmt"
	"github.com/angryscorp/alert-metrics/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

type MetricRouter struct {
	router  *gin.Engine
	storage domain.MetricStorage
}

func NewMetricRouter(router *gin.Engine, storage domain.MetricStorage) *MetricRouter {
	mr := MetricRouter{router: router, storage: storage}
	mr.registerPing()
	mr.registerGetMetric()
	mr.registerGetAllMetrics()
	mr.registerUpdateMetrics()
	mr.registerFetchMetricsJSON()
	mr.registerUpdateMetricsJSON()
	return &mr
}

func (mr *MetricRouter) Run(addr string) (err error) {
	return mr.router.Run(addr)
}

func (mr *MetricRouter) registerPing() {
	mr.router.GET("/ping", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
}

func (mr *MetricRouter) registerGetMetric() {
	mr.router.GET("/value/:metricType/:metricName", func(c *gin.Context) {
		metricType, err := domain.NewMetricType(c.Param("metricType"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		metrics, ok := mr.storage.GetMetrics(metricType, c.Param("metricName"))
		if !ok {
			c.Status(http.StatusNotFound)
			return
		}

		c.String(http.StatusOK, metrics.StringValue())
	})
}

func (mr *MetricRouter) registerGetAllMetrics() {
	mr.router.GET("/", func(c *gin.Context) {
		htmlContent := "<h3>Current metrics</h3><ul>"
		for k, v := range mr.storage.GetAllMetrics() {
			htmlContent += fmt.Sprintf("<li>%s: %s</li>", k, v)
		}
		htmlContent += "</ul>"
		c.Data(http.StatusOK, "text/html", []byte(htmlContent))
	})
}

func (mr *MetricRouter) registerUpdateMetrics() {
	mr.router.POST("/update/:metricType/:metricName/:metricValue", func(c *gin.Context) {
		if err := mr.update(c.Param("metricType"), c.Param("metricName"), c.Param("metricValue")); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		c.Status(http.StatusOK)
	})
}

func (mr *MetricRouter) registerFetchMetricsJSON() {
	mr.router.POST("/value/", func(c *gin.Context) {
		if err := mr.verifyContentTypeIsJSON(c.Request); err != nil {
			c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": err.Error()})
			return
		}

		var metrics domain.Metrics
		if err := c.ShouldBindJSON(&metrics); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		res, ok := mr.storage.GetMetrics(metrics.MType, metrics.ID)
		if ok {
			c.JSON(http.StatusOK, res)
		}
		c.Header("Content-Type", "application/json")
		c.Status(http.StatusNotFound)
	})
}

func (mr *MetricRouter) registerUpdateMetricsJSON() {
	mr.router.POST("/update/", func(c *gin.Context) {
		if err := mr.verifyContentTypeIsJSON(c.Request); err != nil {
			c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": err.Error()})
			return
		}

		var metrics domain.Metrics
		if err := c.ShouldBindJSON(&metrics); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		metrics, err := mr.updateMetrics(metrics)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, metrics)
	})
}

func (mr *MetricRouter) update(rawMetricType string, metricName string, metricValue string) error {
	metrics, err := domain.NewMetrics(rawMetricType, metricName, metricValue)
	if err != nil {
		return err
	}

	return mr.storage.UpdateMetrics(*metrics)
}

func (mr *MetricRouter) updateMetrics(metrics domain.Metrics) (domain.Metrics, error) {
	err := mr.storage.UpdateMetrics(metrics)
	if err != nil {
		return domain.Metrics{}, err
	}

	res, ok := mr.storage.GetMetrics(metrics.MType, metrics.ID)
	if !ok {
		return domain.Metrics{}, errors.New("failed to get updated metrics")
	}
	return res, nil
}

func (mr *MetricRouter) verifyContentTypeIsJSON(r *http.Request) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return errors.New("Content-Type must be application/json")
	}
	return nil
}
