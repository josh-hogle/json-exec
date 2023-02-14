package config

// VersionConfig contains the options for the "version" command.
type VersionConfig struct {
	// Plaintext indicates whether or not to display the version information in plaintext.
	Plaintext bool `yaml:"plaintext"`

	// Verbose indicates whether or not build and release information should be printed along with the version.
	Verbose bool `yaml:"verbose"`
}
