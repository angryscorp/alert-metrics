package serverconfig

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
)

type ServerConfig struct {
	Address string `env:"ADDRESS"`
}

func ParseConfig() (ServerConfig, error) {
	address := flag.String("a", "localhost:8080", "HTTP server address (default: localhost:8080)")

	// Parsing
	flag.Parse()

	// Unknown flags
	if len(flag.Args()) > 0 {
		return ServerConfig{}, fmt.Errorf("unknown flag or argument %s", flag.Args())
	}

	config := ServerConfig{Address: *address}

	// ENV vars
	err := env.Parse(&config)
	if err != nil {
		return ServerConfig{}, err
	}

	return config, nil
}
