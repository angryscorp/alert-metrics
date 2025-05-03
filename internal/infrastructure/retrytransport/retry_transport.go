package retrytransport

import (
	"bytes"
	"github.com/rs/zerolog"
	"io"
	"net/http"
	"time"
)

type RetryTransport struct {
	transport      http.RoundTripper
	retryIntervals []time.Duration
	logger         zerolog.Logger
}

func New(
	transport http.RoundTripper,
	retryIntervals []time.Duration,
	logger zerolog.Logger,
) *RetryTransport {
	if transport == nil {
		transport = http.DefaultTransport
	}

	return &RetryTransport{
		transport:      transport,
		retryIntervals: append([]time.Duration{0}, retryIntervals...),
		logger:         logger,
	}
}

func (rt *RetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	var reqBody []byte
	if req.Body != nil {
		reqBody, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		_ = req.Body.Close()
	}

	for n, interval := range rt.retryIntervals {
		if reqBody != nil {
			req.Body = io.NopCloser(bytes.NewBuffer(reqBody))
		}

		if n > 0 {
			rt.logger.Warn().Msgf("retrying request %d in %f", n, interval.Seconds())
			time.Sleep(interval)
		}

		resp, err = rt.transport.RoundTrip(req)
		if err == nil {
			return resp, nil
		}

		rt.logger.Error().Err(err).Msgf("failed to send request %s", err)
	}

	return nil, err
}
