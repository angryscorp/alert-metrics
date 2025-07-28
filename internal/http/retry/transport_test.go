package retry

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("with custom transport", func(t *testing.T) {
		transport := &mockTransport{}
		retryIntervals := []time.Duration{time.Second}

		result := New(transport, retryIntervals, zerolog.Nop())

		require.NotNil(t, result)
		assert.Equal(t, transport, result.transport)
		assert.Equal(t, []time.Duration{0, time.Second}, result.retryIntervals)
	})
}

func TestTransport_RoundTrip(t *testing.T) {
	tests := []struct {
		name          string
		mockResponses []mockResponse
		expectedCalls int
		expectError   bool
	}{
		{
			name: "success on first try",
			mockResponses: []mockResponse{
				{resp: &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader("ok"))}, err: nil},
			},
			expectedCalls: 1,
			expectError:   false,
		},
		{
			name: "success after retry",
			mockResponses: []mockResponse{
				{resp: nil, err: errors.New("network error")},
				{resp: &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader("ok"))}, err: nil},
			},
			expectedCalls: 2,
			expectError:   false,
		},
		{
			name: "all retries fail",
			mockResponses: []mockResponse{
				{resp: nil, err: errors.New("error 1")},
				{resp: nil, err: errors.New("error 2")},
			},
			expectedCalls: 2,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTransport := &mockTransport{responses: tt.mockResponses}
			transport := New(mockTransport, []time.Duration{time.Millisecond}, zerolog.Nop())

			req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
			resp, err := transport.RoundTrip(req)
			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(resp.Body)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, resp)
			}

			assert.Equal(t, tt.expectedCalls, mockTransport.callCount)
		})
	}
}

type mockResponse struct {
	resp *http.Response
	err  error
}

type mockTransport struct {
	responses []mockResponse
	callCount int
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.callCount >= len(m.responses) {
		return nil, errors.New("no more responses")
	}

	response := m.responses[m.callCount]
	m.callCount++
	return response.resp, response.err
}
