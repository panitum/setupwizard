package custom

import (
	"fmt"
	customPkg "setupwizard/internal/custom"

	"github.com/spf13/cobra"
)

var CustomCmd = &cobra.Command{
	Use:   "custom",
	Short: "Executing custom .sh-scripts",
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeCustomCommands(); err != nil {
			panic(err)
		}
	},
}

func executeCustomCommands() error {
	scripts, err := customPkg.NewCustom()
	if err != nil {
		return err
	}

	if scripts.Count() < 1 {
		fmt.Println("No custom scripts found.")
		return nil
	}

	if err := scripts.Execute(); err != nil {
		return err
	}

	return nil
}
