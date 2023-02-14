package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/josh-hogle/zerolog/v2/log"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// Load simply loads settings from the config file.
//
// The file must be a valid YAML-formatted file and must have a .yaml or .yml extension.
//
// The following order of precedence is used in determining the path to the config file:
//   - the value passed to this function if it is not empty
//   - the value of the prefixed CONFIG_FILE environment variable (this is the EnvPrefix followed by an underscore
//     followed by CONFIG_FILE)
//   - the default configuration file
//
// The following order of precedence is used to determine the location of the default config file:
//   - the current working directory
//   - the default configuration directory inside the user's home directory
//
// If a configuration file is specified explicitly and that file does not exist, this function will return an error.
// If one is not specified and the default configuration file cannot be found, no error will be returned. Instead
// default values will be used for the configuration settings.
func Load(file string) error {
	// set the path to the config file
	var configFile string
	usingDefaultConfig := false
	configFileEnv := os.Getenv(fmt.Sprintf("%s_CONFIG_FILE", EnvPrefix))
	if file != "" {
		configFile, err := filepath.Abs(os.ExpandEnv(file))
		if err != nil {
			return fmt.Errorf("failed to determine absolute path of config file '%s': %s", file, err.Error())
		}
		_viper.SetConfigFile(configFile)
	} else if configFileEnv != "" {
		configFile, err := filepath.Abs(os.ExpandEnv(configFileEnv))
		if err != nil {
			return fmt.Errorf("failed to determine absolute path of config file '%s': %s", configFileEnv, err.Error())
		}
		_viper.SetConfigFile(configFile)
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("unable to get current working directory: %s", err.Error())
		}
		user, err := user.Current()
		if err != nil {
			return fmt.Errorf("unable to retrieve current user home directory: %s", err.Error())
		}
		_viper.AddConfigPath(cwd)
		_viper.AddConfigPath(filepath.Join(user.HomeDir, DefaultConfigFolder))
		_viper.SetConfigName(DefaultConfigName)
		usingDefaultConfig = true
	}

	// read the configuration file
	if err := _viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if !usingDefaultConfig {
				return fmt.Errorf("failed to load config file '%s': file not found", configFile)
			}
		} else {
			return fmt.Errorf("failed to load config file '%s': %s", _viper.ConfigFileUsed(), err.Error())
		}
	}

	// watch for configuration file changes
	//
	// Note that we don't use the built-in WatchConfig() function because it will spit out a text error message
	// if there is a parse error in the YAML, which could break the output if it's expected to be JSON.
	watchConfig()
	return Unmarshal()
}

// Unmarshal simply unmarshals the viper configuration into the global app object.
func Unmarshal() error {
	// viper uses a mapstructure to attempt cover multiple file formats, environment variables
	// and command-line options. mapstructure decoding can be a bit flaky so we rely on more
	// tried and true methods of unmarshalling YAML by first converting the mapstructure back
	// to YAML and then unmarshalling it, which also gives us the added benefit of utilizing
	// UnmarshalYAML() functions for our custom objects.
	configYaml, err := yaml.Marshal(_viper.AllSettings())
	if err != nil {
		return fmt.Errorf("failed to parse configuration settings: %s", err.Error())
	}

	var config AppConfig
	if err := yaml.Unmarshal(configYaml, &config); err != nil {
		return fmt.Errorf("failed to parse configuration settings: %s", err.Error())
	}
	_config = &config

	// now that we've parsed the settings, we can start using logger to output mesages
	log.Debug().
		Interface("config_settings", _config).
		Msgf("Successfully parsed configuration")
	return nil
}

// reloadConfig handles reloading the configuration settings after the config file changes.
func reloadConfig() {
	// BUG: odd behavior here but fsnotify seems to trigger 2 write events when the config is edited
	//      so Unmarshal will get called twice, but it shouldn't be an issue since we are locking the
	//      Unmarshal with a mutex so that we don't get concurrency issues.
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	if err := Unmarshal(); err != nil {
		log.
			Error().Stack().
			Err(err).
			Msgf("aborting reload of configuration file settings: %s", err.Error())
		log.
			Warn().
			Msg("rolling back changes to original settings")
	}
}

// watchConfig watches the config file for changes.
//
// This function is taken almost directly from the viper source code. It is used in place of the WatchConfig()
// function of viper because that function may write output that is not properly in compliance with the output
// format configured for the application. This function fixes that issue.
func watchConfig() {
	initWG := sync.WaitGroup{}
	initWG.Add(1)
	go func() {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Error().Err(err).Msgf("failed to create watcher: %s", err.Error())
			log.Fatal()
		}
		defer watcher.Close()

		// we have to watch the entire directory to pick up renames/atomic saves in a cross-platform way
		filename := _viper.ConfigFileUsed()
		configFile := filepath.Clean(filename)
		configDir, _ := filepath.Split(configFile)
		realConfigFile, _ := filepath.EvalSymlinks(filename)

		eventsWG := sync.WaitGroup{}
		eventsWG.Add(1)
		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok { // 'Events' channel is closed
						eventsWG.Done()
						return
					}
					currentConfigFile, _ := filepath.EvalSymlinks(filename)
					// we only care about the config file with the following cases:
					// 1 - if the config file was modified or created
					// 2 - if the real path to the config file changed (eg: k8s ConfigMap replacement)
					const writeOrCreateMask = fsnotify.Write | fsnotify.Create
					if (filepath.Clean(event.Name) == configFile &&
						event.Op&writeOrCreateMask != 0) ||
						(currentConfigFile != "" && currentConfigFile != realConfigFile) {
						realConfigFile = currentConfigFile
						err := _viper.ReadInConfig()
						if err != nil {
							log.
								Error().Stack().
								Err(err).
								Msgf("aborting reload of configuration file settings: %s", err.Error())
							log.
								Warn().
								Msg("no changes have been made to the configuration")
						} else {
							reloadConfig()
						}
					} else if filepath.Clean(event.Name) == configFile &&
						event.Op&fsnotify.Remove&fsnotify.Remove != 0 {
						eventsWG.Done()
						return
					}

				case err, ok := <-watcher.Errors:
					if ok { // 'Errors' channel is not closed
						log.
							Error().Stack().
							Err(err).
							Msgf("error occurred watching for file changes: %s", err.Error())
					}
					eventsWG.Done()
					return
				}
			}
		}()
		watcher.Add(configDir)
		initWG.Done()   // done initializing the watch in this go routine, so the parent routine can move on...
		eventsWG.Wait() // now, wait for event loop to end in this go-routine...
	}()
	initWG.Wait() // make sure that the go routine above fully ended before returning
}
