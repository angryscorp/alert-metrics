package main

import (
	"flag"
	"fmt"
	"github.com/angryscorp/alert-metrics/internal/infrastructure/httplogger"
	"github.com/angryscorp/alert-metrics/internal/infrastructure/metricrouter"
	"github.com/angryscorp/alert-metrics/internal/infrastructure/metricstorage"
	"github.com/angryscorp/alert-metrics/internal/infrastructure/serverconfig"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"os"
)

func main() {
	flags, err := serverconfig.ParseConfig()
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err.Error())
		flag.Usage()
		os.Exit(1)
	}

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	router := gin.New()
	router.
		Use(httplogger.New(logger)).
		Use(gin.Recovery()).
		Use(gzip.Gzip(gzip.DefaultCompression))

	var mr = metricrouter.NewMetricRouter(
		router,
		metricstorage.NewMemStorage(),
	)

	err = mr.Run(flags.Address)
	if err != nil {
		panic(err)
	}
}
