package version

import (
	"fmt"

	"github.com/josh-hogle/json-exec/internal/app"
	"github.com/josh-hogle/json-exec/internal/config"
	"github.com/josh-hogle/json-exec/internal/errors"
	"github.com/spf13/cobra"
	"go.innotegrity.dev/zerolog/log"
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
			Use:   "version",
			Short: "Display application version information",
			Long: `
This command displays basic or detailed version information about the application.`,
		},

		main: main,
	}
	cmd.Run = cmd.run

	// flags managed via viper
	viper := config.Viper()
	flags := cmd.Flags()

	flags.BoolP("verbose", "v", false, "display full version information including build and release date")
	viper.SetDefault("version.verbose", false)
	viper.BindPFlag("version.verbose", flags.Lookup("verbose"))

	flags.BoolP("plaintext", "p", false, "output plaintext instead of JSON")
	viper.SetDefault("version.plaintext", false)
	viper.BindPFlag("version.plaintext", flags.Lookup("plaintext"))

	return cmd
}

// run simply executes the command.
func (c *Command) run(cmd *cobra.Command, args []string) {
	cfg := config.Get()

	logger := log.With().
		Str("version", app.Version).
		Logger()
	output := fmt.Sprintf("%s version %s", app.Title, app.Version)

	if cfg.Version.Verbose {
		logger = log.With().
			Str("build", app.Build).
			Str("release_date", app.ReleaseDate).
			Str("developer_build", app.IsDevelopment).
			Logger()

		output += fmt.Sprintf(" build %s", app.Build)
		if app.ReleaseDate == "Unreleased" {
			output += " (Unreleased)"
		} else {
			output += fmt.Sprintf(" (Released %s)", app.ReleaseDate)
		}
		if app.IsDevBuild {
			output += " [developer build]"
		}
	}

	if cfg.Version.Plaintext {
		fmt.Printf("%s\n", output)
	} else {
		logger.Info().Msg(output)
	}

	c.main.SetExitCode(errors.None)
}
