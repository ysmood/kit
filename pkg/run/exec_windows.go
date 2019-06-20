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
	buf := make([]byte, 32*1024)
	prefixBuf := []byte(prefix)
	bufOut := make([]byte, 32*1024+len(prefixBuf))

	for {
		n, err := reader.Read(buf)
		if err == io.EOF {
			break
		}

		bufOutIndex := 0
		bufOutIndex += copy(bufOut[bufOutIndex:], prefixBuf)
		for _, r := range string(buf[:n]) {
			if err != nil {
				_, _ = gos.Stdout.Write(buf[:n])
				return err
			}
			if r == '\n' {
				bufOutIndex += copy(bufOut[bufOutIndex:], prefixBuf)
			}
			bufOutIndex += copy(bufOut[bufOutIndex:], []byte(string(r)))
		}
		_, _ = gos.Stdout.Write(bufOut[:bufOutIndex])
	}

	return nil
}

// KillTree kill process and all its children process
func KillTree(pid int) error {
	return exec.Command("taskkill", "/t", "/f", "/pid", strconv.Itoa(pid)).Run()
}
