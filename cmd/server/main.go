package main

import (
	"github.com/angryscorp/alert-metrics/internal/infrastructure/metricrouter"
	"github.com/angryscorp/alert-metrics/internal/infrastructure/metricstorage"
	"github.com/gin-gonic/gin"
)

func main() {
	var mr = metricrouter.NewMetricRouter(
		gin.Default(),
		metricstorage.NewMemStorage(),
	)

	err := mr.Run(":8080")
	if err != nil {
		panic(err)
	}
}
