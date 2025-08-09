package server

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
			"ADDRESS":           "example.com:9090",
			"STORE_INTERVAL":    "600",
			"FILE_STORAGE_PATH": "/tmp/metrics.dump",
			"RESTORE":           "true",
			"DATABASE_DSN":      "postgres://user:pass@localhost/db",
			"KEY":               "secret123",
			"CRYPTO_KEY":        "file.pem",
		}

		expected := Config{
			Address:                "example.com:9090",
			StoreIntervalInSeconds: 600,
			FileStoragePath:        "/tmp/metrics.dump",
			ShouldRestore:          true,
			DatabaseDSN:            "postgres://user:pass@localhost/db",
			HashKey:                "secret123",
			PathToCryptoKey:        "file.pem",
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
