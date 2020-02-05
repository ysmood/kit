package utils

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/k0kubun/pp"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"github.com/mgutz/ansi"
)

var goos = runtime.GOOS

// Stdout ...
var Stdout = stdout()

// Stderr ...
var Stderr = stderr()

// Dump ...
func Dump(val interface{}) {
	E(pp.Println(val))
}

// Sdump ...
func Sdump(val interface{}) string {
	return pp.Sprint(val)
}

// Log log to stdout with timestamp
func Log(v ...interface{}) {
	t := time.Now().Format("[2006-01-02 15:04:05]")
	v = append([]interface{}{C(t, "7")}, v...)
	E(fmt.Fprintln(Stdout, v...))
}

// Err log to stderr with timestamp and stack trace
func Err(v ...interface{}) {
	t := time.Now().Format("[2006-01-02 15:04:05]")
	v = append(v, "\n"+string(debug.Stack()))
	v = append([]interface{}{C(t, "7")}, v...)

	E(fmt.Fprintln(Stderr, v...))
}

// ClearScreen ...
func ClearScreen() error {
	if goos == "windows" {
		_, err := os.Stdout.WriteString("\n\n\n\n\n")
		return err
	}

	print("\033[H\033[2J")
	return nil
}

// C color terminal string
func C(str interface{}, color string) string {
	return ansi.Color(fmt.Sprint(str), color)
}

func stdout() io.Writer {
	if goos == "windows" {
		fd := os.Stdout.Fd()
		if !isatty.IsTerminal(fd) && !isatty.IsCygwinTerminal(fd) {
			return colorable.NewNonColorable(os.Stdout)
		}
	}
	return colorable.NewColorableStdout()
}

func stderr() io.Writer {
	if goos == "windows" {
		fd := os.Stderr.Fd()
		if !isatty.IsTerminal(fd) && !isatty.IsCygwinTerminal(fd) {
			return colorable.NewNonColorable(os.Stderr)
		}
	}
	return colorable.NewColorableStderr()
}
