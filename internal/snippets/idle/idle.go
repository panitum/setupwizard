package idle

import (
	"errors"
	"os/exec"
	"runtime"
)

const ErrUnsupportedOs = "unsupported os"

func IncreaseIdleTimeout() error {
	cmd, err := setIdleTimeout("180")
	if err != nil {
		return err
	}

	return cmd.Run()
}

func SetDefaultIdleTimeout() error {
	cmd, err := setIdleTimeout("30")
	if err != nil {
		return err
	}

	return cmd.Run()
}

func setIdleTimeout(minutes string) (*exec.Cmd, error) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("sudo", "pmset", "-a", "displaysleep", minutes)
	default:
		return nil, errors.New(ErrUnsupportedOs)
	}

	return cmd, nil
}
