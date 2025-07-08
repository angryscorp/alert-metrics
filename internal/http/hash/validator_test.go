package hash

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewHashValidator(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		hashKey        string
		requestBody    string
		hashHeader     string
		expectedStatus int
		shouldCallNext bool
	}{
		{
			name:           "empty hash key - should skip validation",
			hashKey:        "",
			requestBody:    "test body",
			hashHeader:     "invalid_hash",
			expectedStatus: http.StatusOK,
			shouldCallNext: true,
		},
		{
			name:           "missing hash header - should skip validation",
			hashKey:        "secret",
			requestBody:    "test body",
			hashHeader:     "",
			expectedStatus: http.StatusOK,
			shouldCallNext: true,
		},
		{
			name:           "valid hash - should pass validation",
			hashKey:        "secret",
			requestBody:    "test body",
			hashHeader:     computeHash("test body", "secret"),
			expectedStatus: http.StatusOK,
			shouldCallNext: true,
		},
		{
			name:           "invalid hash - should fail validation",
			hashKey:        "secret",
			requestBody:    "test body",
			hashHeader:     "invalid_hash",
			expectedStatus: http.StatusBadRequest,
			shouldCallNext: false,
		},
		{
			name:           "empty body with valid hash",
			hashKey:        "secret",
			requestBody:    "",
			hashHeader:     computeHash("", "secret"),
			expectedStatus: http.StatusOK,
			shouldCallNext: true,
		},
		{
			name:           "different hash key - should fail validation",
			hashKey:        "secret",
			requestBody:    "test body",
			hashHeader:     computeHash("test body", "different_secret"),
			expectedStatus: http.StatusBadRequest,
			shouldCallNext: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			nextCalled := false

			validator := NewHashValidator(tt.hashKey)

			router.POST("/test", validator, func(c *gin.Context) {
				nextCalled = true
				c.Status(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(tt.requestBody))
			if tt.hashHeader != "" {
				req.Header.Set("HashSHA256", tt.hashHeader)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.shouldCallNext, nextCalled)
		})
	}
}

// Helper function to compute expected hash
func computeHash(body, key string) string {
	h := sha256.New()
	h.Write([]byte(body))
	h.Write([]byte(key))
	return fmt.Sprintf("%x", h.Sum(nil))
}
