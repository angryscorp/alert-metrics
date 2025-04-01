package serverflags

import (
	"flag"
	"fmt"
)

type ServerFlags struct {
	Address string
}

func SetupAndParseFlags() (ServerFlags, error) {
	// Flags
	address := flag.String("a", "localhost:8080", "HTTP server address (default: localhost:8080)")

	// Parsing
	flag.Parse()

	// Unknown flags
	if len(flag.Args()) > 0 {
		return ServerFlags{}, fmt.Errorf("unknown flag or argument %s", flag.Args())
	}

	return ServerFlags{Address: *address}, nil
}
