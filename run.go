// +build !windows

package gokit

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/kr/pty"
)

func run(prefix string, cmd *exec.Cmd) error {
	p, err := pty.Start(cmd)
	if err != nil {
		return err
	}

	// Make sure to close the pty at the end.
	defer func() { p.Close() }() // Best effort.

	// Handle pty size.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, p); err != nil {
				Log("[exec] resize pty:", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH // Initial resize.

	go func() {
		io.Copy(p, os.Stdin)
	}()

	reader := bufio.NewReader(p)
	newline := true
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			os.Stdout.Write([]byte(string(r)))
			break
		}
		if newline {
			os.Stdout.Write([]byte(prefix))
			newline = false
		}
		if r == '\n' {
			newline = true
		}
		os.Stdout.Write([]byte(string(r)))
	}

	return cmd.Wait()
}
