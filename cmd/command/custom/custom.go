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
		executeCustomCommands()
	},
}

func executeCustomCommands() {
	scripts, err := customPkg.NewCustom()
	if err != nil {
		fmt.Println(err)
		return
	}

	if scripts.Count() < 1 {
		fmt.Println("No custom scripts found.")
		return
	}

	for _, script := range scripts.GetScripts() {
		err := customPkg.RunScript(script)
		if err != nil {
			fmt.Printf("\n[Ошибка выполнения sh-скрипта]: %v\n", err)
			return
		}
	}
}
