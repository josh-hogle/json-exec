package config

import "go.innotegrity.dev/zerolog"

const (
	// DefaultConfigFolder is the name of the config folder in the user's home directory.
	DefaultConfigFolder = ".json-exec"

	// DefaultConfigName is the default configuration file name without an extension.
	DefaultConfigName = "json-exec"

	// DefaultLogLevel is the default logging level.
	DefaultLogLevel = zerolog.InfoLevel

	// DefaultLogLevelFieldName is the name of the level field in log messages.
	DefaultLogLevelFieldName = "@level"

	// DefaultLogMessageFieldName is the name of the message field in log messages.
	DefaultLogMessageFieldName = "@message"

	// DefaultLogTimestampFieldName is the name of the timestamp field in log messages.
	DefaultLogTimestampFieldName = "@timestamp"

	// EnvPrefix is the prefix used for configuration via environment variables.
	EnvPrefix = "JSON_EXEC"
)
