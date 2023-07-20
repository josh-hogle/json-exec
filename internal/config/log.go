package config

import (
	"strings"

	"go.innotegrity.dev/zerolog/log"
)

// defaultLogWriter is just a wrapper so we can catch any messages from internal libraries that use the
// standard Go log interface and output them through zerolog.
type defaultLogWriter struct{}

// Write just writes the string into a new log message on "info" level.
func (w defaultLogWriter) Write(p []byte) (n int, err error) {
	message := strings.TrimSuffix(string(p), "\n")
	log.Info().
		Str("log_message", message).
		Msgf("internal Go log message detected: %s", message)
	return len(p), nil
}
