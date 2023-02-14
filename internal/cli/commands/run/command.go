package run

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/josh-hogle/json-exec/internal/app"
	"github.com/josh-hogle/json-exec/internal/config"
	"github.com/josh-hogle/json-exec/internal/errors"
	"github.com/josh-hogle/zerolog/v2/log"
	"github.com/spf13/cobra"
)

// Command is the object for executing the actual command
type Command struct {
	*cobra.Command

	// unexported members
	main app.Main
}

// NewCommand creates a new Command object.
func NewCommand(main app.Main) *Command {
	if main == nil { // should never happen
		panic("null 'main' pointer passed to NewCommand()")
	}
	cmd := &Command{
		Command: &cobra.Command{
			Use:   "run [flags] <command> [command args]",
			Short: "Executes an arbitrary system command with optional flags",
			Long: `
run will execute the given system command with any flags passed to the command. Be sure
to use -- before the system command when it requires its own set of flags.`,
			Args: cobra.MinimumNArgs(1),
		},

		main: main,
	}
	cmd.Run = cmd.run

	// flags managed via viper
	viper := config.Viper()
	flags := cmd.Flags()

	flags.Bool("ignore-stdout", false, "ignore stdout output from the command")
	viper.SetDefault("run.ignore_stdout", false)
	viper.BindPFlag("run.ignore_stdout", flags.Lookup("ignore-stdout"))

	flags.Bool("ignore-stderr", false, "ignore stderr output from the command")
	viper.SetDefault("run.ignore_stderr", false)
	viper.BindPFlag("run.ignore_stderr", flags.Lookup("ignore-stderr"))

	return cmd
}

// run simply executes the command.
func (c *Command) run(cmd *cobra.Command, args []string) {
	restoreLogger, _ := log.ReplaceGlobal(
		log.With().
			Str("command", args[0]).
			Interface("args", args[1:]).
			Logger(),
	)
	defer restoreLogger()
	cfg := config.Get()

	// run the command
	var stdout, stderr bytes.Buffer
	command := exec.Command(args[0], args[1:]...)
	if cfg.Run.IgnoreStdout {
		command.Stdout = nil
	} else {
		command.Stdout = &stdout
	}
	if cfg.Run.IgnoreStderr {
		command.Stderr = nil
	} else {
		command.Stderr = &stderr
	}
	errorMessage := ""
	exitCode := errors.None
	log.Info().Msgf("executing command: %s", strings.Join(args, " "))
	if err := command.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
			errorMessage = exitErr.Error()
		} else {
			exitCode = 99
			errorMessage = fmt.Sprintf("%s", err)
		}
	}

	// print the results
	logger := log.With().
		Int("exit_code", exitCode).
		Logger()
	if errorMessage != "" {
		logger = log.With().
			Str("error_message", errorMessage).
			Logger()
	}
	if !cfg.Run.IgnoreStdout {
		logger = log.With().
			Str("stdout", stdout.String()).
			Logger()
	}
	if !cfg.Run.IgnoreStderr {
		logger = log.With().
			Str("stderr", stderr.String()).
			Logger()
	}

	if exitCode != errors.None {
		logger.Warn().Msgf("command exited with non-zero exit code %d", exitCode)
	} else {
		logger.Info().Msg("command completed successfully")
	}
	c.main.SetExitCode(exitCode)
}
