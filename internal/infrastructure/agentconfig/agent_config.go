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
}

func New() (AgentConfig, error) {
	// Flags
	address := flag.String("a", "localhost:8080", "HTTP agent address (default: localhost:8080)")
	pollIntervalInSeconds := flag.Int("p", 2, "Poll interval in seconds (default: 2)")
	reportIntervalInSeconds := flag.Int("r", 10, "Report interval in seconds (default: 10)")

	// Parsing
	flag.Parse()

	// Unknown flags
	if len(flag.Args()) > 0 {
		return AgentConfig{}, fmt.Errorf("unknown flag or argument %s", flag.Args())
	}

	config := AgentConfig{
		Address:                 *address,
		ReportIntervalInSeconds: *reportIntervalInSeconds,
		PollIntervalInSeconds:   *pollIntervalInSeconds,
	}

	// ENV vars
	err := env.Parse(&config)
	if err != nil {
		return AgentConfig{}, err
	}

	return config, nil
}
