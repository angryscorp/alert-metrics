package agent

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	t.Run("environment variables override", func(t *testing.T) {
		// Reset flag package
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		envVars := map[string]string{
			"ADDRESS":         "example.com:9090",
			"POLL_INTERVAL":   "5",
			"REPORT_INTERVAL": "30",
			"KEY":             "secret123",
			"RATE_LIMIT":      "20",
			"CRYPTO_KEY":      "file.pem",
		}

		expected := Config{
			Address:                 "example.com:9090",
			PollIntervalInSeconds:   5,
			ReportIntervalInSeconds: 30,
			HashKey:                 "secret123",
			RateLimit:               20,
			PathToCryptoKey:         "file.pem",
		}

		for key, value := range envVars {
			_ = os.Setenv(key, value)
		}

		oldArgs := os.Args
		os.Args = []string{"test"}
		defer func() { os.Args = oldArgs }()

		config, err := NewConfig()

		require.NoError(t, err)
		assert.Equal(t, expected, config)
	})
}
