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
				Err("error resizing pty: %s", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH // Initial resize.

	go func() { io.Copy(p, os.Stdin) }()

	scanner := bufio.NewScanner(p)
	for scanner.Scan() {
		os.Stdout.Write([]byte(prefix + scanner.Text() + "\n"))
	}

	return cmd.Wait()
}
