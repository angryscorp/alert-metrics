package router

import "github.com/gin-gonic/gin"

type PingHandler interface {
	Ping(c *gin.Context)
}

type MetricsHandler interface {
	GetMetric(c *gin.Context)
	GetAllMetrics(c *gin.Context)
	UpdateMetrics(c *gin.Context)
}

type MetricsJSONHandler interface {
	FetchMetricsJSON(c *gin.Context)
	UpdateMetricsJSON(c *gin.Context)
	BatchUpdateFetchMetrics(c *gin.Context)
}
