package config

import (
	"fmt"
	"strings"

	"go.innotegrity.dev/zerolog"
	"go.innotegrity.dev/zerolog/log"
	"gopkg.in/yaml.v3"
)

// GlobalConfig contains the global options for the application.
type GlobalConfig struct {
	// ExtraFields is a mapping of extra fields to add to the output messages.
	ExtraFields map[string]string `yaml:"extra_fields"`

	// LevelFieldName contains the name of the level field.
	LevelFieldName string `yaml:"level_field_name"`

	// LogLevel contains the actual zerolog level for logging output messages.
	LogLevel zerolog.Level

	// LogLevelRaw represents the string version of the logging level.
	LogLevelRaw string `yaml:"log_level"`

	// MessageFieldName contains the name of the message field.
	MessageFieldName string `yaml:"message_field_name"`

	// TimestampFieldName contains the name of the timestamp field.
	TimestampFieldName string `yaml:"timestamp_field_name"`
}
type _yamlGlobalConfig GlobalConfig // wrapper to avoid infinite recursion

// UnmarshalYAML decodes the raw YAML into the object.
//
// It converts any raw values to their corresponding actual values and then performs validation on the
// object member values. It may set default values as well, if necessary.
func (c *GlobalConfig) UnmarshalYAML(value *yaml.Node) error {
	// unmarshal into a temporary object so we don't overwrite existing settins if the operation fails
	var cfg _yamlGlobalConfig
	if err := value.Decode(&cfg); err != nil {
		return err
	}
	*c = GlobalConfig(cfg)

	// set log level
	if strings.EqualFold(c.LogLevelRaw, "none") {
		log.SetLevel(zerolog.Disabled)
	} else {
		level, err := zerolog.ParseLevel(c.LogLevelRaw)
		if err != nil {
			return fmt.Errorf("failed to parse log level '%s': %s", c.LogLevelRaw, err.Error())
		}
		c.LogLevel = level
		log.SetLevel(level)
	}

	// set field names
	if c.LevelFieldName == "" {
		zerolog.LevelFieldName = DefaultLogLevelFieldName
	} else {
		zerolog.LevelFieldName = c.LevelFieldName
	}
	if c.MessageFieldName == "" {
		zerolog.MessageFieldName = DefaultLogMessageFieldName
	} else {
		zerolog.MessageFieldName = c.MessageFieldName
	}
	if c.TimestampFieldName == "" {
		zerolog.TimestampFieldName = DefaultLogTimestampFieldName
	} else {
		zerolog.TimestampFieldName = c.TimestampFieldName
	}

	// add extra fields to the output
	logger := log.Logger
	if c.ExtraFields != nil {
		for k, v := range c.ExtraFields {
			logger = log.With().Str(k, v).Logger()
		}
		log.ReplaceGlobal(logger)
	}
	return nil
}
