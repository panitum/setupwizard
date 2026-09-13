package command

import (
	"os"
	"setupwizard/cmd/command/config"
	"setupwizard/cmd/command/custom"
	"setupwizard/cmd/command/install"
	cfg "setupwizard/internal/config"
	customPkg "setupwizard/internal/custom"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "setupwizard",
	Short: "Let's set up your computer.",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	if err := customPkg.Init(); err != nil {
		panic(err)
	}

	rootCmd.AddCommand(config.ConfigCmd)

	if cfg.IsExist() {
		rootCmd.AddCommand(install.InstallCmd)
	}

	rootCmd.AddCommand(custom.CustomCmd)
}
