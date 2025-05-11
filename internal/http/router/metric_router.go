package router

import (
	"context"
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

func New(
	router *gin.Engine,
	storage domain.MetricStorage,
) *MetricRouter {
	mr := MetricRouter{router: router, storage: storage}

	mr.registerNoRoutes()
	mr.registerPing()
	mr.registerGetMetric()
	mr.registerGetAllMetrics()
	mr.registerUpdateMetrics()
	mr.registerFetchMetricsJSON()
	mr.registerUpdateMetricsJSON()
	mr.registerBatchUpdateFetchMetrics()

	return &mr
}

func (mr *MetricRouter) Run(addr string) (err error) {
	return mr.router.Run(addr)
}

func (mr *MetricRouter) registerNoRoutes() {
	mr.router.NoRoute(func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})

	// CI/CD bug? should be handled by NoRoute, but it's not
	mr.router.Any("/updater/*any", func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})

	mr.router.POST("/update/counter/", func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})

	mr.router.POST("/update/gauge/", func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})
}

func (mr *MetricRouter) registerPing() {
	mr.router.GET("/ping", func(c *gin.Context) {
		if err := mr.storage.Ping(c.Request.Context()); err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
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

		metrics, ok := mr.storage.GetMetric(c.Request.Context(), metricType, c.Param("metricName"))
		if !ok {
			c.Status(http.StatusNotFound)
			return
		}

		c.String(http.StatusOK, metrics.StringValue())
	})
}

func (mr *MetricRouter) registerGetAllMetrics() {
	mr.router.GET("/", func(c *gin.Context) {
		allMetrics := mr.storage.GetAllMetrics(c.Request.Context())
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
	})
}

func (mr *MetricRouter) registerUpdateMetrics() {
	mr.router.POST("/update/:metricType/:metricName/:metricValue", func(c *gin.Context) {
		if err := mr.update(c.Request.Context(), c.Param("metricType"), c.Param("metricName"), c.Param("metricValue")); err != nil {
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

		var metrics domain.Metric
		if err := c.ShouldBindJSON(&metrics); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		res, ok := mr.storage.GetMetric(c.Request.Context(), metrics.MType, metrics.ID)
		if ok {
			c.JSON(http.StatusOK, res)
			return
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

		var metric domain.Metric
		if err := c.ShouldBindJSON(&metric); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		metric, err := mr.updateMetrics(c.Request.Context(), metric)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, metric)
	})
}

func (mr *MetricRouter) registerBatchUpdateFetchMetrics() {
	mr.router.POST("/updates/", func(c *gin.Context) {
		if err := mr.verifyContentTypeIsJSON(c.Request); err != nil {
			c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": err.Error()})
			return
		}

		var metrics []domain.Metric
		if err := c.ShouldBindJSON(&metrics); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := mr.storage.UpdateMetrics(c.Request.Context(), metrics)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, metrics)
	})
}

func (mr *MetricRouter) update(ctx context.Context, rawMetricType string, metricName string, metricValue string) error {
	metrics, err := domain.NewMetrics(rawMetricType, metricName, metricValue)
	if err != nil {
		return err
	}

	return mr.storage.UpdateMetric(ctx, *metrics)
}

func (mr *MetricRouter) updateMetrics(ctx context.Context, metrics domain.Metric) (domain.Metric, error) {
	err := mr.storage.UpdateMetric(ctx, metrics)
	if err != nil {
		return domain.Metric{}, err
	}

	res, ok := mr.storage.GetMetric(ctx, metrics.MType, metrics.ID)
	if !ok {
		return domain.Metric{}, errors.New("failed to get updated metrics")
	}
	return res, nil
}

func (mr *MetricRouter) verifyContentTypeIsJSON(r *http.Request) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return errors.New("Content-Type must be application/json")
	}

	return nil
}
