package agentconfig

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
)

type AgentConfig struct {
	Address                 string `env:"ADDRESS"`
	PollIntervalInSeconds   int    `env:"POLL_INTERVAL"`
	ReportIntervalInSeconds int    `env:"REPORT_INTERVAL"`
	HashKey                 string `env:"KEY"`
}

func New() (AgentConfig, error) {
	address := flag.String("a", "localhost:8080", "HTTP agent address (default: localhost:8080)")
	pollIntervalInSeconds := flag.Int("p", 2, "Poll interval in seconds (default: 2)")
	reportIntervalInSeconds := flag.Int("r", 10, "Report interval in seconds (default: 10)")
	hashKey := flag.String("k", "", "Key for calculating hash (default: none)")

	flag.Parse()

	// Unknown flags
	if len(flag.Args()) > 0 {
		return AgentConfig{}, fmt.Errorf("unknown flag or argument %s", flag.Args())
	}

	config := AgentConfig{
		Address:                 *address,
		ReportIntervalInSeconds: *reportIntervalInSeconds,
		PollIntervalInSeconds:   *pollIntervalInSeconds,
		HashKey:                 *hashKey,
	}

	// ENV vars
	err := env.Parse(&config)
	if err != nil {
		return AgentConfig{}, err
	}

	return config, nil
}
