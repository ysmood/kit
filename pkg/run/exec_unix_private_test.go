// +build !windows

package run

import (
	"errors"
	"io"
	"os"
	"testing"

	gos "github.com/ysmood/gokit/pkg/os"
	"golang.org/x/crypto/ssh/terminal"
)

func TestRestoreState(t *testing.T) {
	restoreState(&terminal.State{})
}

type testWriter struct {
	err error
}

func (t testWriter) Write(p []byte) (n int, err error) {
	return 0, t.err
}

func (t testWriter) Read(p []byte) (n int, err error) {
	return 0, t.err
}

func TestStdinPiper(t *testing.T) {
	stdinWriter = testWriter{err: errors.New("err")}
	old := os.Stdin
	os.Stdin, _ = os.Open(gos.ThisFilePath())
	defer func() { os.Stdin = old }()
	stdinPiper()
}

func TestPipeToStdoutWithPrefixReadEOF(t *testing.T) {
	pipeToStdoutWithPrefix("", testWriter{err: io.EOF})
}

func TestPipeToStdoutWithPrefixReadErr(t *testing.T) {
	pipeToStdoutWithPrefix("", testWriter{err: errors.New("err")})
}
