package buildinfo

import "fmt"

type BuildInfo struct {
	version string
	date    string
	commit  string
}

func New(version, date, commit string) BuildInfo {
	return BuildInfo{
		version: version,
		date:    date,
		commit:  commit,
	}
}

func (buildInfo BuildInfo) String() string {
	if buildInfo.version == "" {
		buildInfo.version = "N/A"
	}

	if buildInfo.date == "" {
		buildInfo.date = "N/A"
	}

	if buildInfo.commit == "" {
		buildInfo.commit = "N/A"
	}

	return fmt.Sprintf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildInfo.version, buildInfo.date, buildInfo.commit)
}
