package logger

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		responseStatus int
		setupHandler   func(c *gin.Context)
	}{
		{
			name:           "GET request with 200 status",
			method:         "GET",
			path:           "/test",
			responseStatus: http.StatusOK,
			setupHandler: func(c *gin.Context) {
				c.String(http.StatusOK, "test response")
			},
		},
		{
			name:           "PUT request with 400 status",
			method:         "PUT",
			path:           "/api/users/123",
			responseStatus: http.StatusBadRequest,
			setupHandler: func(c *gin.Context) {
				c.String(http.StatusBadRequest, "bad request")
			},
		},
		{
			name:           "DELETE request with 404 status",
			method:         "DELETE",
			path:           "/api/users/456",
			responseStatus: http.StatusNotFound,
			setupHandler: func(c *gin.Context) {
				c.Status(http.StatusNotFound)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var logBuffer bytes.Buffer
			logger := zerolog.New(&logBuffer).With().Timestamp().Logger()

			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.Use(New(logger))
			router.Handle(tt.method, tt.path, tt.setupHandler)

			req, err := http.NewRequest(tt.method, tt.path, nil)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			logOutput := logBuffer.String()
			assert.Contains(t, logOutput, `"method":"`+tt.method+`"`)
			assert.Contains(t, logOutput, `"uri":"`+tt.path+`"`)
			assert.Contains(t, logOutput, `"duration":`)
			assert.Contains(t, logOutput, `"status":`+strconv.Itoa(tt.responseStatus))
			assert.Contains(t, logOutput, `"size":`)
			assert.Contains(t, logOutput, `"level":"info"`)
		})
	}
}
