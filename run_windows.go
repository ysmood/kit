// +build windows

package gokit

import (
	"bufio"
	"io"
	"os"
	"os/exec"
)

// The pty lib doesn't support Windows, so we just pipe everything
func run(prefix string, cmd *exec.Cmd) error {
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

	scanner := bufio.NewScanner(io.MultiReader(stderr, stdout))
	for scanner.Scan() {
		os.Stdout.Write([]byte(prefix + scanner.Text() + "\n"))
	}
	return nil
}
