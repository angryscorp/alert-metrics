package crypto

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/angryscorp/alert-metrics/internal/domain"
)

func DecrypterMiddleware(decrypter domain.Decrypter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Checking if the request is encrypted
		if c.GetHeader("Content-Encoding") != "encrypted" {
			c.Next()
			return
		}

		// Reading encrypted data
		encryptedData, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
			return
		}
		_ = c.Request.Body.Close()

		// Decrypting
		decryptedData, err := decrypter.Decrypt(encryptedData)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to decrypt request body"})
			return
		}

		// Replacing the request body with decrypted data
		c.Request.Body = io.NopCloser(bytes.NewReader(decryptedData))
		c.Request.ContentLength = int64(len(decryptedData))

		// Removing the encryption header
		c.Request.Header.Del("Content-Encoding")

		c.Next()
	}
}
