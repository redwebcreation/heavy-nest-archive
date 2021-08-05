package cli

import (
	"github.com/redwebcreation/hez/internal"
	"github.com/redwebcreation/hez/internal/ui"
	"github.com/spf13/cobra"
	"os"
)

var overwriteConfig bool

func runInitCommand(_ *cobra.Command, _ []string) error {
	err := internal.ElevateProcess()
	_, err = os.Stat(internal.ConfigFile)

	if !os.IsNotExist(err) && !overwriteConfig {
		ui.Warning("A config already exists at " + internal.ConfigFile)
		ui.Warning("You may overwrite using the `--force` flag.")
		return nil
	}

	err = os.MkdirAll(internal.ConfigDirectory, 0777)
	if err != nil {
		return err
	}

	err = os.WriteFile(internal.ConfigFile, []byte(`#applications: [ ]
proxy:
  http:
    enabled: true
    port: 80
  https:
    enabled: true
    port: 443
    self_signed: false
  logs:
    level: 0
    redirections:
      - stdout
`), 0777)
	if err != nil {
		return err
	}
	ui.Success("The configuration file has been written at " + internal.ConfigFile)

	return nil
}

func InitCommand() *cobra.Command {
	return internal.CreateCommand(
		&cobra.Command{
			Use:   "init",
			Short: "Creates a default configuration",
		}, func(command *cobra.Command) {
			command.Flags().BoolVarP(&overwriteConfig, "overwriteConfig", "f", false, "Can overwrite a previous config.l")
		}, runInitCommand)
}
