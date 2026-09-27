package idle

import (
	"errors"
	"os/exec"
	"runtime"
)

const ErrUnsupportedOs = "unsupported os"

func IncreaseIdleTimeout() error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("sudo", "pmset", "-a", "displaysleep", "180")
	default:
		return errors.New(ErrUnsupportedOs)
	}

	return cmd.Run()
}
