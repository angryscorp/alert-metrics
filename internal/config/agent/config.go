package agent

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Address                 string `env:"ADDRESS" json:"address"`
	PollIntervalInSeconds   int    `env:"POLL_INTERVAL" json:"poll_interval"`
	ReportIntervalInSeconds int    `env:"REPORT_INTERVAL" json:"report_interval"`
	HashKey                 string `env:"KEY"`
	RateLimit               int    `env:"RATE_LIMIT"`
	PathToCryptoKey         string `env:"CRYPTO_KEY" json:"crypto_key"`
}

func NewConfig() (Config, error) {
	configPath := flag.String("c", "", "Path to the config file")
	flag.StringVar(configPath, "config", "", "Path to the config file")

	address := flag.String("a", "localhost:8080", "HTTP agent address (default: localhost:8080)")
	pollIntervalInSeconds := flag.Int("p", 2, "Poll interval in seconds (default: 2)")
	reportIntervalInSeconds := flag.Int("r", 10, "Report interval in seconds (default: 10)")
	hashKey := flag.String("k", "", "Key for calculating hash (default: none)")
	rateLimit := flag.Int("l", 10, "Rate limit (default: 10)")
	pathToCryptoKey := flag.String("crypto-key", "", "Path to a file with a public key (default: none)")

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
	if *pollIntervalInSeconds != -1 {
		config.PollIntervalInSeconds = *pollIntervalInSeconds
	}
	if *reportIntervalInSeconds != -1 {
		config.ReportIntervalInSeconds = *reportIntervalInSeconds
	}
	if *rateLimit != -1 {
		config.RateLimit = *rateLimit
	}
	if *hashKey != "" {
		config.HashKey = *hashKey
	}
	if *pathToCryptoKey != "" {
		config.PathToCryptoKey = *pathToCryptoKey
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
