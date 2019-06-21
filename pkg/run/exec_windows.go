// +build windows

package run

import (
	"io"
	"os"
	"os/exec"
	"strconv"
)

// The pty lib doesn't support Windows, so we just pipe everything
func run(prefix string, isRaw bool, cmd *exec.Cmd) error {
	cmd.Stdin = os.Stdin

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	pipeToStdoutWithPrefix(prefix, io.MultiReader(stderr, stdout))

	return nil
}

// KillTree kill process and all its children process
func KillTree(pid int) error {
	return exec.Command("taskkill", "/t", "/f", "/pid", strconv.Itoa(pid)).Run()
}
