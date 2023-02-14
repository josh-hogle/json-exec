package cli

import (
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/josh-hogle/json-exec/internal/app"
	"github.com/josh-hogle/json-exec/internal/cli/commands/run"
	"github.com/josh-hogle/json-exec/internal/cli/commands/version"
	"github.com/josh-hogle/json-exec/internal/config"
	"github.com/josh-hogle/json-exec/internal/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// RootCommand is the root command of the application.
type RootCommand struct {
	*cobra.Command

	// unexported members
	configFile string
	exitCode   int
}

// NewRootCommand creates and initializes a new RootCommand object.
func NewRootCommand() *RootCommand {
	cmd := &RootCommand{
		Command: &cobra.Command{
			Use:   fmt.Sprintf("%s [flags]", app.CommandName),
			Short: "Encodes output from a command into a JSON message containing common fields",
			Long: fmt.Sprintf(`
%s will execute a command and redirect its output from both stderr and stdout
and produce 2 JSON messages for a log file: the first containing details about the command
and the second containing the results of the command.`, app.Title),
		},

		exitCode: errors.Usage, // since we have not yet parsed the command-line
	}
	cmd.PersistentPreRunE = cmd.persistentPreRunE
	cmd.SilenceErrors = true

	// flags not managed by viper
	pflags := cmd.PersistentFlags()
	cmd.PersistentFlags().StringVarP(&cmd.configFile, "config-file", "c", "",
		"Path to the configuration settings file")

	// flags managed via viper
	viper := config.Viper()

	pflags.StringP("log-level", "l", config.DefaultLogLevel.String(),
		"adjust output log level - must be one of: debug, info, warn, error, fatal, panic or none")
	viper.SetDefault("global.log_level", config.DefaultLogLevel.String())
	viper.BindPFlag("global.log_level", pflags.Lookup("log-level"))

	pflags.StringToStringP("field", "f", nil,
		"one or more additional fields to include in the output")
	viper.SetDefault("global.extra_fields", nil)
	viper.BindPFlag("global.extra_fields", pflags.Lookup("field"))

	pflags.String("level-field", config.DefaultLogLevelFieldName,
		"alternate name for the level field")
	viper.SetDefault("global.level_field_name", nil)
	viper.BindPFlag("global.level_field_name", pflags.Lookup("level-field"))

	pflags.String("message-field", config.DefaultLogMessageFieldName,
		"alternate name for the message field")
	viper.SetDefault("global.message_field_name", nil)
	viper.BindPFlag("global.message_field_name", pflags.Lookup("message-field"))

	pflags.String("timestamp-field", config.DefaultLogTimestampFieldName,
		"alternate name for the timestamp field")
	viper.SetDefault("global.timestamp_field_name", nil)
	viper.BindPFlag("global.timestamp_field_name", pflags.Lookup("timestamp-field"))

	// add commands
	cmd.AddCommand(
		run.NewCommand(cmd).Command,
		version.NewCommand(cmd).Command,
	)
	return cmd
}

// GetCommandLineFlags returns a map of all global and command-specified flags passed via the command line.
//
// If a flag was not passed, it will not be present in the map. This will allow the caller to quickly
// determine whether or not a flag was specified on the command line rather than having to traverse an
// entire array.
//
// The boolean value the flag is mapped to will always be true.
func (c *RootCommand) GetCommandLineFlags(cmd *cobra.Command) map[string]bool {
	visitedFlags := map[string]bool{}
	cmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		if f.Changed {
			visitedFlags[f.Name] = true
		}
	})
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if f.Changed {
			visitedFlags[f.Name] = true
		}
	})
	return visitedFlags
}

// GetExitCode retrieves the exit code for the application.
func (c *RootCommand) GetExitCode() int {
	return c.exitCode
}

// SetExitCode sets the exit code for the application.
func (c *RootCommand) SetExitCode(code int) {
	c.exitCode = code
}

// persistentPreRunE always runs before any command is executed in order to perform global initialization.
func (c *RootCommand) persistentPreRunE(cmd *cobra.Command, args []string) error {
	// command-line flags were successfully processed so start with a clean exit code
	c.exitCode = errors.None

	// set special environment variables
	env := map[string]string{
		"BUILD":          app.Build,
		"COMMAND_NAME":   app.CommandName,
		"IS_DEVELOPMENT": strconv.FormatBool(app.IsDevBuild),
		"RELEASE_DATE":   app.ReleaseDate,
		"VERSION":        app.Version,
		"GOOS":           runtime.GOOS,
		"GOARCH":         runtime.GOARCH,
	}
	for k, v := range env {
		if err := os.Setenv(fmt.Sprintf("%s_%s", config.EnvPrefix, k), v); err != nil {
			c.exitCode = errors.GeneralFailure
			return fmt.Errorf("failed to set environment variable '%s': %s", k, err.Error())
		}
	}

	// load any settings from config file
	visitedFlags := c.GetCommandLineFlags(c.Command)
	configFile := ""
	if _, ok := visitedFlags["config-file"]; ok {
		configFile = c.configFile
	}
	if err := config.Load(configFile); err != nil {
		c.exitCode = errors.ConfigLoadFailure
		return err
	}
	return nil
}
