package zipper

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
)

type GzipTransport struct {
	Transport http.RoundTripper
}

var _ http.RoundTripper = (*GzipTransport)(nil)

func NewGzipTransport(transport http.RoundTripper) *GzipTransport {
	if transport == nil {
		transport = http.DefaultTransport
	}
	return &GzipTransport{Transport: transport}
}

func (gt *GzipTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Body == nil {
		return gt.Transport.RoundTrip(req)
	}

	contentType := req.Header.Get("Content-Type")
	if !isContentZippable(contentType) {
		return gt.Transport.RoundTrip(req)
	}

	body, err := io.ReadAll(req.Body)
	_ = req.Body.Close()
	if err != nil {
		return nil, err
	}

	if len(body) < minSize {
		req.Body = io.NopCloser(bytes.NewReader(body))
		return gt.Transport.RoundTrip(req)
	}

	compressed, err := gt.zip(body)
	if err != nil {
		return nil, err
	}

	req.Body = io.NopCloser(bytes.NewReader(compressed))
	req.ContentLength = int64(len(compressed))
	req.Header.Set("Content-Encoding", "gzip")

	return gt.Transport.RoundTrip(req)
}

func (gt *GzipTransport) zip(body []byte) ([]byte, error) {
	var compressed bytes.Buffer
	gzipWriter := gzip.NewWriter(&compressed)

	_, err := gzipWriter.Write(body)
	if err != nil {
		return nil, err
	}

	if err = gzipWriter.Close(); err != nil {
		return nil, err
	}

	return compressed.Bytes(), nil
}
