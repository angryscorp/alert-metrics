package gzipper

import (
	"compress/gzip"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UnzipMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("Content-Encoding") == "gzip" {
			reader, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
			c.Request.Body = io.NopCloser(reader)
			c.Request.Header.Del("Content-Encoding")
			c.Request.Header.Del("Content-Length")

		}
		c.Next()
	}
}
