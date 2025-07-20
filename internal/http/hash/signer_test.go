package hash

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHashTransport(t *testing.T) {
	transport := &mockTransport{}
	hashKey := "secret"

	signer := NewHashTransport(transport, hashKey)

	assert.NotNil(t, signer)
	assert.Equal(t, transport, signer.transport)
	assert.Equal(t, hashKey, signer.hashKey)
}

func TestSigner_RoundTrip(t *testing.T) {
	tests := []struct {
		name           string
		hashKey        string
		requestBody    string
		hasBody        bool
		expectedHeader bool
		expectError    bool
	}{
		{
			name:           "empty hash key - should skip signing",
			hashKey:        "",
			requestBody:    "test body",
			hasBody:        true,
			expectedHeader: false,
			expectError:    false,
		},
		{
			name:           "valid request with body - should add hash header",
			hashKey:        "secret",
			requestBody:    "test body",
			hasBody:        true,
			expectedHeader: true,
			expectError:    false,
		},
		{
			name:           "empty body - should add hash header",
			hashKey:        "secret",
			requestBody:    "",
			hasBody:        true,
			expectedHeader: true,
			expectError:    false,
		},
		{
			name:           "large body - should add hash header",
			hashKey:        "secret",
			requestBody:    strings.Repeat("a", 1024),
			hasBody:        true,
			expectedHeader: true,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTransport := &mockTransport{
				roundTripFunc: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(strings.NewReader("response")),
					}, nil
				},
			}

			signer := NewHashTransport(mockTransport, tt.hashKey)

			req := httptest.NewRequest(http.MethodPost, "http://example.com", nil)
			if tt.hasBody {
				req.Body = io.NopCloser(strings.NewReader(tt.requestBody))
			}

			resp, err := signer.RoundTrip(req)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, resp)

			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(resp.Body)

			hashHeader := req.Header.Get("HashSHA256")
			if tt.expectedHeader {
				assert.NotEmpty(t, hashHeader)
			} else {
				assert.Empty(t, hashHeader)
			}
		})
	}
}

// mockTransport is a mock implementation of http.RoundTripper
type mockTransport struct {
	roundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.roundTripFunc != nil {
		return m.roundTripFunc(req)
	}
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader("default response")),
	}, nil
}
