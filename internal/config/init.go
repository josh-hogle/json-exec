package config

import (
	golog "log"
	"strings"

	"github.com/spf13/viper"
)

var (
	_config *AppConfig
	_viper  *viper.Viper
)

func init() {
	// configure viper settings
	viper.SupportedExts = []string{"yaml", "yml"}
	_viper = viper.New()
	_viper.SetConfigType("yaml")
	_viper.SetEnvPrefix(EnvPrefix)
	_viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	_viper.AutomaticEnv()

	// re-route the default Go logger to our logger
	var w defaultLogWriter
	l := golog.Default()
	l.SetFlags(0)
	l.SetPrefix("")
	l.SetOutput(w)

	// create the global configuration object
	_config = &AppConfig{}
}

// Get returns the one and only AppConfig object.
func Get() *AppConfig {
	return _config
}

// Viper returns the one and only Viper object.
func Viper() *viper.Viper {
	return _viper
}
