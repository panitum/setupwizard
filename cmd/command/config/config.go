package config

import (
	"github.com/spf13/cobra"
)

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Config actions",
}

func init() {
	ConfigCmd.AddCommand(generateCmd)
}
