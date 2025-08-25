package server

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Address                string `env:"ADDRESS" json:"address"`
	StoreIntervalInSeconds int    `env:"STORE_INTERVAL" json:"store_interval"`
	FileStoragePath        string `env:"FILE_STORAGE_PATH" json:"store_file"`
	ShouldRestore          bool   `env:"RESTORE" json:"restore"`
	DatabaseDSN            string `env:"DATABASE_DSN" json:"database_dsn"`
	HashKey                string `env:"KEY"`
	PathToCryptoKey        string `env:"CRYPTO_KEY" json:"crypto_key"`
	TrustedSubnet          string `env:"TRUSTED_SUBNET" json:"trusted_subnet"`
	UseGRPC                bool   `env:"USE_GRPC" json:"use_grpc"`
	GRPCAddress            string `env:"GRPC _ADDRESS" json:"grpc_address"`
}

func NewConfig() (Config, error) {
	configPath := flag.String("c", "", "Path to the config file")
	flag.StringVar(configPath, "config", "", "Path to the config file")

	address := flag.String("a", "localhost:8080", "HTTP server address (default: localhost:8080)")
	storeIntervalInSeconds := flag.Int("i", 300, "Store interval in seconds (default: 300)")
	fileStoragePath := flag.String("f", "alert_monitoring_metrics.dump", "File storage path (default: alert_monitoring_metrics.dump)")
	shouldRestore := flag.Bool("r", false, "Restore from file (default: false)")
	databaseDSN := flag.String("d", "", "Database DSN (default: empty, file storage will be used)")
	hashKey := flag.String("k", "", "Key for calculating hash (default: none)")
	pathToCryptoKey := flag.String("crypto-key", "", "Path to a file with a private key (default: none)")
	isSubnetTrusted := flag.String("t", "", "Path to a file with a public key (default: none)")
	useGRPC := flag.Bool("g", false, "Use also GRPC for incoming requests (default: false)")
	grpcAddress := flag.String("ga", "localhost:443", "gRPC server address (default: localhost:443)")

	flag.Parse()

	// Unknown flags
	if len(flag.Args()) > 0 {
		return Config{}, fmt.Errorf("unknown flag or argument %s", flag.Args())
	}

	configFilePath := *configPath
	if configFilePath == "" {
		configFilePath = os.Getenv("CONFIG")
	}

	config := Config{}

	// Config file
	err := config.loadFromFile(configFilePath)
	if err != nil {
		return Config{}, fmt.Errorf("failed to load config from file: %w", err)
	}

	// CLI flags
	if *address != "" {
		config.Address = *address
	}

	if *storeIntervalInSeconds != -1 {
		config.StoreIntervalInSeconds = *storeIntervalInSeconds
	}

	if *fileStoragePath != "" {
		config.FileStoragePath = *fileStoragePath
	}

	if flag.Lookup("r").Value.String() == "true" {
		config.ShouldRestore = *shouldRestore
	}

	if *databaseDSN != "" {
		config.DatabaseDSN = *databaseDSN
	}

	if *hashKey != "" {
		config.HashKey = *hashKey
	}

	if *pathToCryptoKey != "" {
		config.PathToCryptoKey = *pathToCryptoKey
	}

	if *isSubnetTrusted != "" {
		config.TrustedSubnet = *isSubnetTrusted
	}

	if flag.Lookup("g").Value.String() == "true" {
		config.UseGRPC = *useGRPC
	}

	if *grpcAddress != "" {
		config.GRPCAddress = *grpcAddress
	}

	// ENV vars
	err = env.Parse(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func (cfg *Config) loadFromFile(filePath string) error {
	if filePath == "" {
		return nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("failed to parse JSON config: %w", err)
	}

	return nil
}
