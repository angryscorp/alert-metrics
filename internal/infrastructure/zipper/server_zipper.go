package zipper

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func NewHandlerZipper() gin.HandlerFunc {
	zipper := gzip.Gzip(gzip.DefaultCompression,
		gzip.WithCustomShouldCompressFn(func(c *gin.Context) bool {
			contentType := c.Writer.Header().Get("Content-Type")
			return isContentZippable(contentType) && c.Writer.Size() >= minSize
		}),
	)

	return zipper
}
