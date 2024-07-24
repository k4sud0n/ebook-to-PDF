package capture

import (
	"os/exec"
)

func MacOS(fileName string, rect string) error {
	cmd := exec.Command("screencapture", "-R", rect, fileName)
	return cmd.Run()
}
