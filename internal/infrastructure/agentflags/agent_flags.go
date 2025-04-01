package agentflags

import (
	"flag"
	"fmt"
)

type AgentFlags struct {
	Address                 string
	PollIntervalInSeconds   int
	ReportIntervalInSeconds int
}

func SetupAndParseFlags() (AgentFlags, error) {
	// Flags
	address := flag.String("a", "localhost:8080", "HTTP agent address (default: localhost:8080)")
	pollIntervalInSeconds := flag.Int("p", 2, "Poll interval in seconds (default: 2)")
	reportIntervalInSeconds := flag.Int("r", 10, "Report interval in seconds (default: 10)")

	// Parsing
	flag.Parse()

	// Unknown flags
	if len(flag.Args()) > 0 {
		return AgentFlags{}, fmt.Errorf("unknown flag or argument %s", flag.Args())
	}

	return AgentFlags{
		Address:                 *address,
		ReportIntervalInSeconds: *reportIntervalInSeconds,
		PollIntervalInSeconds:   *pollIntervalInSeconds,
	}, nil
}
