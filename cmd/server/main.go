package main

import (
	"flag"
	"fmt"
	"github.com/angryscorp/alert-metrics/internal/infrastructure/metricrouter"
	"github.com/angryscorp/alert-metrics/internal/infrastructure/metricstorage"
	"github.com/angryscorp/alert-metrics/internal/infrastructure/serverconfig"
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	flags, err := serverconfig.ParseConfig()
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err.Error())
		flag.Usage()
		os.Exit(1)
	}

	var mr = metricrouter.NewMetricRouter(
		gin.Default(),
		metricstorage.NewMemStorage(),
	)

	err = mr.Run(flags.Address)
	if err != nil {
		panic(err)
	}
}
