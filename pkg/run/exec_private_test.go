package run

import (
	"errors"
	"os"
	"testing"

	gos "github.com/ysmood/gokit/pkg/os"
	"golang.org/x/crypto/ssh/terminal"
)

func TestRestoreState(t *testing.T) {
	restoreState(&terminal.State{})
}

type testWriter struct{}

func (t testWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("err")
}

func TestStdinPiper(t *testing.T) {
	stdinWriter = testWriter{}
	old := os.Stdin
	os.Stdin, _ = os.Open(gos.ThisFilePath())
	defer func() { os.Stdin = old }()
	stdinPiper()
}
