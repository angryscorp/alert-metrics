package logger

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"time"
)

func New(logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logger.Info().
			Str("method", c.Request.Method).
			Str("uri", c.Request.URL.Path).
			Dur("duration", time.Since(start)).
			Int("status", c.Writer.Status()).
			Int("size", c.Writer.Size()).
			Msg("")
	}
}
