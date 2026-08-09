package command

import (
	"fmt"
	"os"
	"setupwizard/cmd/command/config"
	"setupwizard/cmd/command/install"
	cfg "setupwizard/internal/config"
	"strconv"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "setupwizard",
	Short: "Let's set up your computer.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cfg.IsExist() {
			fmt.Println("Config exists.")
		} else {
			fmt.Println("Config file not exists!!!")
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(config.GetCmd())

	if cfg.IsExist() {
		rootCmd.AddCommand(install.GetCmd())
	}

	generateLong()
}

func generateLong() {
	for i, command := range rootCmd.Commands() {
		rootCmd.Long += strconv.Itoa(i+1) + ". " + command.Use + " — " + command.Short + "\n"
	}
}
