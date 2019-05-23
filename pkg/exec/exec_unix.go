// +build !windows

package exec

import (
	"bufio"
	"io"
	"os"
	os_exec "os/exec"
	"os/signal"
	"sync"
	"syscall"

	"github.com/kr/pty"
	gos "github.com/ysmood/gokit/pkg/os"
	"golang.org/x/crypto/ssh/terminal"
)

var rawLock = sync.Mutex{}

func run(prefix string, isRaw bool, cmd *os_exec.Cmd) error {
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
		for {
			if _, ok := <-ch; !ok {
				return
			}
			_ = pty.InheritSize(os.Stdin, p)
		}
	}()
	ch <- syscall.SIGWINCH // Initial resize.

	if isRaw {
		// Set stdin in raw mode.
		rawLock.Lock()
		oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			gos.Log("[exec] set stdin to raw mode:", err)
		}
		defer func() {
			if oldState != nil {
				_ = terminal.Restore(int(os.Stdin.Fd()), oldState)
			}
			rawLock.Unlock()
		}() // Best effort.
	}

	go func() {
		_, _ = io.Copy(p, os.Stdin)
	}()

	reader := bufio.NewReader(p)
	newline := true
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			_, _ = gos.Stdout.Write([]byte(string(r)))
			break
		}
		if newline {
			_, _ = gos.Stdout.Write([]byte(prefix))
			newline = false
		}
		if r == '\n' {
			newline = true
		}
		_, _ = gos.Stdout.Write([]byte(string(r)))
	}

	signal.Stop(ch)
	close(ch)

	return cmd.Wait()
}

// KillTree kill process and all its children process
func KillTree(pid int) error {
	group, err := os.FindProcess(-1 * pid)
	if err != nil {
		return err
	}

	return group.Signal(syscall.SIGTERM)
}
