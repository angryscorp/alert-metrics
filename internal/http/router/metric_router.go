package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type MetricRouter struct {
	engine *gin.Engine
}

func New(
	engine *gin.Engine,
) *MetricRouter {
	mr := MetricRouter{engine: engine}
	mr.registerNoRoutes()
	return &mr
}

func (mr *MetricRouter) Run(addr string) (err error) {
	return mr.engine.Run(addr)
}

func (mr *MetricRouter) RegisterPingHandler(handler PingHandler) {
	mr.engine.GET("/ping", handler.Ping)
}

func (mr *MetricRouter) RegisterMetricsHandler(handler MetricsHandler) {
	mr.engine.GET("/value/:metricType/:metricName", handler.GetMetric)
	mr.engine.GET("/", handler.GetAllMetrics)
	mr.engine.POST("/update/:metricType/:metricName/:metricValue", handler.UpdateMetrics)
	mr.engine.POST("/value/", handler.FetchMetricsJSON)
	mr.engine.POST("/update/", handler.UpdateMetricsJSON)
	mr.engine.POST("/updates/", handler.BatchUpdateFetchMetrics)
}

func (mr *MetricRouter) registerNoRoutes() {
	mr.engine.NoRoute(func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})

	// CI/CD bug? should be handled by NoRoute, but it's not
	mr.engine.Any("/updater/*any", func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})

	mr.engine.POST("/update/counter/", func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})

	mr.engine.POST("/update/gauge/", func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})
}
