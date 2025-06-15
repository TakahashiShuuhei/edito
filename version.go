package main

import "fmt"

// Version information for edito
// These values can be overridden at build time using -ldflags
var (
	Version   = "0.3.0"
	BuildDate = "unknown"
	GitCommit = "unknown"
)

// GetVersion returns the version string
func GetVersion() string {
	if GitCommit != "unknown" && GitCommit != "" && len(GitCommit) >= 7 {
		return Version + "-" + GitCommit[:7]
	}
	return Version
}

// GetVersionInfo returns detailed version information
func GetVersionInfo() string {
	return fmt.Sprintf("edito version %s\nBuild date: %s\nGit commit: %s", 
		GetVersion(), BuildDate, GitCommit)
}