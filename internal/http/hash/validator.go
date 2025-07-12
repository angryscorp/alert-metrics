package hash

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewHashValidator(hashKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if hashKey == "" {
			c.Next()
			return
		}

		receivedHash := c.GetHeader("HashSHA256")
		if receivedHash == "" {
			c.Next()
			return
		}

		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		h := sha256.New()
		h.Write(bodyBytes)
		h.Write([]byte(hashKey))
		computedHash := fmt.Sprintf("%x", h.Sum(nil))

		if receivedHash != computedHash {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		c.Next()
	}
}
