package subnet

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		trustedSubnet  string
		realIP         string
		expectedStatus int
	}{
		{
			name:           "Empty trusted subnet- should pass",
			trustedSubnet:  "",
			realIP:         "192.168.1.1",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "IP in trusted subnet- should pass",
			trustedSubnet:  "192.168.1.0/24",
			realIP:         "192.168.1.1",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "IP not in trusted subnet- should return 403",
			trustedSubnet:  "192.168.1.0/24",
			realIP:         "10.0.0.1",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Missing X-Real-IP header - should return 403",
			trustedSubnet:  "192.168.1.0/24",
			realIP:         "",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Invalid IP in X-Real-IP - should return 403",
			trustedSubnet:  "192.168.1.0/24",
			realIP:         "invalid-ip",
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.Use(NewTrustedSubnetMiddleware(tt.trustedSubnet))
			r.GET("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			req := httptest.NewRequest("GET", "/test", nil)
			if tt.realIP != "" {
				req.Header.Set("X-Real-IP", tt.realIP)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
