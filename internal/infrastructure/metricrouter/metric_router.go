package metricrouter

import (
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

		val, ok := mr.storage.Get(metricType, c.Param("metricName"))
		if !ok {
			c.Status(http.StatusNotFound)
			return
		}
		c.String(http.StatusOK, val)
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

func (mr *MetricRouter) update(rawMetricType string, metricName string, metricValue string) error {
	metricType, err := domain.NewMetricType(rawMetricType)
	if err != nil {
		return err
	}

	return mr.storage.Update(metricType, metricName, metricValue)
}
