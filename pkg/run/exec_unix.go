// +build !windows

package run

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"

	"github.com/kr/pty"
	gos "github.com/ysmood/gokit/pkg/os"
	"golang.org/x/crypto/ssh/terminal"
)

var rawLock = sync.Mutex{}

func run(prefix string, isRaw bool, cmd *exec.Cmd) error {
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
		rawLock.Lock()
		defer rawLock.Unlock()
		// Set stdin in raw mode.
		oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			gos.Log("[exec] set stdin to raw mode:", err)
		}
		// Best effort
		defer restoreState(oldState)
	}

	stdinWriter = p
	if !stdinPiperRunning {
		go stdinPiper()
	}

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

var stdinWriter io.Writer
var stdinPiperRunning = false

// This will cause race condition, but normally we don't want two process handle
// the stdin at the same time.
func stdinPiper() {
	stdinPiperRunning = true
	buf := make([]byte, 1024)
	for {
		nr, er := os.Stdin.Read(buf)
		if nr > 0 {
			nw, ew := stdinWriter.Write(buf[0:nr])
			if ew != nil || nr != nw {
				break
			}
		}
		if er != nil {
			break
		}
	}
	stdinPiperRunning = false
}

func restoreState(oldState *terminal.State) {
	if oldState != nil {
		_ = terminal.Restore(int(os.Stdin.Fd()), oldState)
	}
}

// KillTree kill process and all its children process
func KillTree(pid int) error {
	group, _ := os.FindProcess(-1 * pid)

	return group.Signal(syscall.SIGTERM)
}
