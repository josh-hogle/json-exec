package config

// RunConfig contains the options for the "run" command.
type RunConfig struct {
	// IgnoreStderr indicates whether or not to ignore output from stderr.
	IgnoreStderr bool `yaml:"ignore_stderr"`

	// IgnoreStdout indicates whether or not to ignore output from stdout.
	IgnoreStdout bool `yaml:"ignore_stdout"`
}
