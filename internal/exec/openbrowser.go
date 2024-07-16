package exec

import (
	"fmt"
	"gitlab_tui/internal/logger"
	"os/exec"
	"runtime"
)

func Openbrowser(url string) {
	var err error
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		logger.Error(fmt.Errorf("unsupported platform"))
	}

	err = cmd.Start()
	if err != nil {
		logger.Error(err)
	}
	err = cmd.Wait()
	if err != nil {
		logger.Error(err)
	}
}
