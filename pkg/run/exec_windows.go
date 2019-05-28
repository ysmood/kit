// +build windows

package run

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"strconv"

	gos "github.com/ysmood/gokit/pkg/os"
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

	reader := bufio.NewReader(io.MultiReader(stderr, stdout))
	newline := true
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			gos.Stdout.Write([]byte(string(r)))
			break
		}
		if newline {
			gos.Stdout.Write([]byte(prefix))
			newline = false
		}
		if r == '\n' {
			newline = true
		}
		gos.Stdout.Write([]byte(string(r)))
	}

	return nil
}

// KillTree kill process and all its children process
func KillTree(pid int) error {
	return exec.Command("taskkill", "/t", "/f", "/pid", strconv.Itoa(pid)).Run()
}
