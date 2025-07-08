package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPingHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name            string
		setupMock       func(*MockMetricStorage)
		setupRequest    func() *http.Request
		expectedStatus  int
		validateContext bool
	}{
		{
			name: "successful ping",
			setupMock: func(m *MockMetricStorage) {
				m.On("Ping", mock.AnythingOfType("context.backgroundCtx")).
					Return(nil)
			},
			setupRequest: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
				return req
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "storage timeout error",
			setupMock: func(m *MockMetricStorage) {
				m.On("Ping", mock.AnythingOfType("context.backgroundCtx")).
					Return(context.DeadlineExceeded)
			},
			setupRequest: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
				return req
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "database connection refused",
			setupMock: func(m *MockMetricStorage) {
				m.On("Ping", mock.AnythingOfType("context.backgroundCtx")).
					Return(errors.New("connection refused"))
			},
			setupRequest: func() *http.Request {
				req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
				return req
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockStorage := &MockMetricStorage{}
			tc.setupMock(mockStorage)

			handler := NewPingHandler(mockStorage)

			req := tc.setupRequest()
			w := httptest.NewRecorder()

			router := gin.New()
			router.GET("/ping", handler.Ping)

			// Act
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.expectedStatus, w.Code, "unexpected status code")
			assert.Empty(t, w.Body.String(), "expected empty response body")

			mockStorage.AssertExpectations(t)
		})
	}

	t.Run("constructor", func(t *testing.T) {
		mockStorage := &MockMetricStorage{}
		handler := NewPingHandler(mockStorage)

		assert.NotNil(t, handler)
		assert.Equal(t, mockStorage, handler.storage)
	})

	t.Run("interface compliance", func(t *testing.T) {
		mockStorage := &MockMetricStorage{}
		handler := NewPingHandler(mockStorage)

		assert.Implements(t, (*interface{ Ping(*gin.Context) })(nil), handler)
	})
}
