package install

import (
	"fmt"
	"os"
	"setupwizard/internal/config"
	"setupwizard/internal/http/download"
	"strconv"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb/v8"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Let's set up your computer.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Installation has been started.")
		install()
		fmt.Println("Installation has been completed.")
	},
}

func GetCmd() *cobra.Command {
	return installCmd
}

func init() {
	generateLong()
}

func generateLong() {
	for i, command := range installCmd.Commands() {
		installCmd.Long += strconv.Itoa(i+1) + ". " + command.Use + " — " + command.Short + "\n"
	}
}

func install() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	downloadApps(cfg.Apps)
}

func downloadApps(apps []config.App) {
	p := mpb.New(mpb.WithWidth(60), mpb.WithRefreshRate(180*time.Millisecond))

	var wg sync.WaitGroup
	sem := make(chan struct{}, 3)

	for _, app := range apps {
		wg.Add(1)
		sem <- struct{}{}

		go func() {
			defer wg.Done()
			defer func() { <-sem }()

			download.File(app.Link, p)
		}()
	}

	wg.Wait()
	p.Wait()
}
