package crypto

import (
	"bytes"
	"io"
	"net/http"

	"github.com/angryscorp/alert-metrics/internal/domain"
)

func EncryptorMiddleware(encryptor domain.Encrypter) func(http.RoundTripper) http.RoundTripper {
	return func(next http.RoundTripper) http.RoundTripper {
		return &encryptorTransport{
			encryptor: encryptor,
			next:      next,
		}
	}
}

type encryptorTransport struct {
	encryptor domain.Encrypter
	next      http.RoundTripper
}

func (t *encryptorTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Body == nil {
		return t.next.RoundTrip(req)
	}

	// Reading the original request body
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	_ = req.Body.Close()

	// Encrypting
	encryptedData, err := t.encryptor.Encrypt(body)
	if err != nil {
		return nil, err
	}

	// Replacing the request body with encrypted data
	req.Body = io.NopCloser(bytes.NewReader(encryptedData))
	req.ContentLength = int64(len(encryptedData))

	// Adding the encryption header
	req.Header.Set("Content-Encoding", "encrypted")

	return t.next.RoundTrip(req)
}
