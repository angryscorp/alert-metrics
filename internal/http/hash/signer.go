package hash

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
)

type Signer struct {
	transport http.RoundTripper
	hashKey   string
}

func NewHashTransport(transport http.RoundTripper, hashKey string) *Signer {
	return &Signer{
		transport: transport,
		hashKey:   hashKey,
	}
}

func (t *Signer) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Body == nil || t.hashKey == "" {
		return t.transport.RoundTrip(req)
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}

	req.Body = io.NopCloser(bytes.NewBuffer(body))

	h := sha256.New()
	h.Write(body)
	h.Write([]byte(t.hashKey))

	hashStr := fmt.Sprintf("%x", h.Sum(nil))
	req.Header.Set("HashSHA256", hashStr)

	return t.transport.RoundTrip(req)
}
