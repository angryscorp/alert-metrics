package zipper

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

const minSize = 1024

func NewZipper() gin.HandlerFunc {
	zipper := gzip.Gzip(gzip.DefaultCompression,
		gzip.WithCustomShouldCompressFn(func(c *gin.Context) bool {
			contentType := c.Writer.Header().Get("Content-Type")
			isTextOrJSON := contentType == "application/json" || contentType == "text/html"
			return isTextOrJSON && c.Writer.Size() >= minSize
		}),
	)

	return zipper
}
