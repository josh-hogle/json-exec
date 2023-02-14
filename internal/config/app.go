package config

// AppConfig holds the configuration settings for the application.
type AppConfig struct {
	// Global holds the global configuration settings.
	Global GlobalConfig `yaml:"global"`

	// Run holds the "run" command configuration settings.
	Run RunConfig `yaml:"run"`

	// Version holds the "version" command configuration settings.
	Version VersionConfig `yaml:"version"`
}
