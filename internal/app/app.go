package app

import "github.com/spf13/cobra"

// Main represents an interface for the main application.
type Main interface {
	// GetCommandLineFlags returns a map of all global and command-specified flags passed via the command line.
	//
	// If a flag was not passed, it should not be present in the map. This will allow the caller to quickly
	// determine whether or not a flag was specified on the command line rather than having to traverse an
	// entire array.
	//
	// The boolean value the flag is mapped to should always be true.
	GetCommandLineFlags(cmd *cobra.Command) map[string]bool

	// GetExitCode retrieves the exit code for the application.
	GetExitCode() int

	// SetExitCode sets the exit code for the application.
	SetExitCode(int)
}
