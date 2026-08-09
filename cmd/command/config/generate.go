package config

import (
	"fmt"
	"os"
	"setupwizard/cmd/dialog"
	"setupwizard/internal/config"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"gen"},
	Short:   "Config actions.",
	Run: func(cmd *cobra.Command, args []string) {
		generateConfig()
	},
}

func generateConfig() {
	if config.IsExist() {
		question := "Config file is already exists. Do you want to overwrite it?"
		if isForce := dialog.Confirm(question); !isForce {
			fmt.Println("Overwriting declined by user.")
			os.Exit(0)
		}
	}

	cfg := config.NewConfig()
	cfg.FillDummy()

	if err := cfg.Save(); err == nil {
		fmt.Println("New config generated successfully.")
	} else {
		fmt.Println(err)
	}
}
