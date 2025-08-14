package router

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

const gracefulShutdownTimeout = 5 * time.Second

type MetricRouter struct {
	engine *gin.Engine
	logger *zerolog.Logger
}

func New(
	engine *gin.Engine,
	logger *zerolog.Logger,
) *MetricRouter {
	mr := MetricRouter{engine: engine, logger: logger}
	mr.registerNoRoutes()
	return &mr
}

func (mr *MetricRouter) Run(addr string, shutdownCh <-chan struct{}) (err error) {
	srv := &http.Server{
		Addr:    addr,
		Handler: mr.engine,
	}

	errCh := make(chan error, 1)
	go func() {
		mr.logger.Info().Str("address", addr).Msg("Starting server")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("failed to start server: %w", err)
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-shutdownCh:
		mr.logger.Info().Msg("Received shutdown signal, starting graceful shutdown")

		// Graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			mr.logger.Error().Err(err).Msg("Server forced to shutdown")
			return fmt.Errorf("server shutdown error: %w", err)
		}

		mr.logger.Info().Msg("Server shutdown completed")
		return nil
	}
}

func (mr *MetricRouter) RegisterPingHandler(handler PingHandler) {
	mr.engine.GET("/ping", handler.Ping)
}

func (mr *MetricRouter) RegisterMetricsHandler(handler MetricsHandler) {
	mr.engine.GET("/value/:metricType/:metricName", handler.GetMetric)
	mr.engine.GET("/", handler.GetAllMetrics)
	mr.engine.POST("/update/:metricType/:metricName/:metricValue", handler.UpdateMetrics)
}

func (mr *MetricRouter) RegisterMetricsJSONHandler(handler MetricsJSONHandler) {
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
