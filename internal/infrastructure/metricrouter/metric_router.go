package metricrouter

import (
	"errors"
	"fmt"
	"github.com/angryscorp/alert-metrics/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

type MetricRouter struct {
	router   *gin.Engine
	storage  domain.MetricStorage
	database interface{ Ping() error }
}

func NewMetricRouter(
	router *gin.Engine,
	storage domain.MetricStorage,
	database interface{ Ping() error },
) *MetricRouter {
	mr := MetricRouter{router: router, storage: storage, database: database}

	mr.registerNoRoutes()
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
		if err := mr.database.Ping(); err != nil {
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

		metrics, ok := mr.storage.GetMetric(metricType, c.Param("metricName"))
		if !ok {
			c.Status(http.StatusNotFound)
			return
		}

		c.String(http.StatusOK, metrics.StringValue())
	})
}

func (mr *MetricRouter) registerGetAllMetrics() {
	mr.router.GET("/", func(c *gin.Context) {
		allMetrics := mr.storage.GetAllMetrics()
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

		var metrics domain.Metric
		if err := c.ShouldBindJSON(&metrics); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		res, ok := mr.storage.GetMetric(metrics.MType, metrics.ID)
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

		var metrics domain.Metric
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

	return mr.storage.UpdateMetric(*metrics)
}

func (mr *MetricRouter) updateMetrics(metrics domain.Metric) (domain.Metric, error) {
	err := mr.storage.UpdateMetric(metrics)
	if err != nil {
		return domain.Metric{}, err
	}

	res, ok := mr.storage.GetMetric(metrics.MType, metrics.ID)
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
