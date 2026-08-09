package config

import (
	"strconv"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Config actions.",
}

func GetCmd() *cobra.Command {
	return configCmd
}

func init() {
	configCmd.AddCommand(generateCmd)

	generateLong()
}

func generateLong() {
	for i, command := range configCmd.Commands() {
		configCmd.Long += strconv.Itoa(i+1) + ". " + command.Use + " — " + command.Short + "\n"
	}
}
