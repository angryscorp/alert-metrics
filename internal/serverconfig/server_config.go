package serverconfig

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
)

type ServerConfig struct {
	Address                string `env:"ADDRESS"`
	StoreIntervalInSeconds int    `env:"STORE_INTERVAL"`
	FileStoragePath        string `env:"FILE_STORAGE_PATH"`
	ShouldRestore          bool   `env:"RESTORE"`
	DatabaseDSN            string `env:"DATABASE_DSN"`
	HashKey                string `env:"KEY"`
}

func New() (ServerConfig, error) {
	address := flag.String("a", "localhost:8080", "HTTP server address (default: localhost:8080)")
	storeIntervalInSeconds := flag.Int("i", 300, "Store interval in seconds (default: 300)")
	fileStoragePath := flag.String("f", "alert_monitoring_metrics.dump", "File storage path (default: alert_monitoring_metrics.dump)")
	shouldRestore := flag.Bool("r", false, "Restore from file (default: false)")
	databaseDSN := flag.String("d", "", "Database DSN (default: empty, file storage will be used)")
	hashKey := flag.String("k", "", "Key for calculating hash (default: none)")

	flag.Parse()

	// Unknown flags
	if len(flag.Args()) > 0 {
		return ServerConfig{}, fmt.Errorf("unknown flag or argument %s", flag.Args())
	}

	config := ServerConfig{
		Address:                *address,
		StoreIntervalInSeconds: *storeIntervalInSeconds,
		FileStoragePath:        *fileStoragePath,
		ShouldRestore:          *shouldRestore,
		DatabaseDSN:            *databaseDSN,
		HashKey:                *hashKey,
	}

	// ENV vars
	err := env.Parse(&config)
	if err != nil {
		return ServerConfig{}, err
	}

	return config, nil
}
