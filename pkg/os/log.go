package os

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/k0kubun/go-ansi"
	ansiColor "github.com/mgutz/ansi"
	"github.com/sanity-io/litter"
)

// Stdout ...
var Stdout = ansi.NewAnsiStdout()

// Stderr ...
var Stderr = ansi.NewAnsiStderr()

// Dump ...
var Dump = litter.Dump

// Sdump ...
var Sdump = litter.Sdump

// Log log to stdout with timestamp
func Log(v ...interface{}) {
	t := time.Now().Format("[2006-01-02 15:04:05]")
	v = append([]interface{}{C(t, "7")}, v...)

	fmt.Fprintln(Stdout, v...)
}

// Err log to stderr with timestamp and stack trace
func Err(v ...interface{}) {
	t := time.Now().Format("[2006-01-02 15:04:05]")
	v = append(v, "\n"+string(debug.Stack()))
	v = append([]interface{}{C(t, "7")}, v...)

	fmt.Fprintln(Stderr, v...)
}

// ClearScreen ...
func ClearScreen() error {
	clsCmd := "clear"

	if runtime.GOOS == "windows" {
		clsCmd = "cls"
	}

	cmd := exec.Command(clsCmd)
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

// C color terminal string
func C(str interface{}, color string) string {
	return ansiColor.Color(fmt.Sprint(str), color)
}
