package os

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/k0kubun/go-ansi"
	"github.com/k0kubun/pp"
	ansiColor "github.com/mgutz/ansi"
	"github.com/ysmood/gokit/pkg/utils"
)

// Stdout ...
var Stdout = ansi.NewAnsiStdout()

// Stderr ...
var Stderr = ansi.NewAnsiStderr()

var goos = runtime.GOOS

// Dump ...
func Dump(val interface{}) {
	utils.E(pp.Println(val))
}

// Sdump ...
func Sdump(val interface{}) string {
	return pp.Sprint(val)
}

// Log log to stdout with timestamp
func Log(v ...interface{}) {
	t := time.Now().Format("[2006-01-02 15:04:05]")
	v = append([]interface{}{C(t, "7")}, v...)

	fmt.Fprintln(Stdout, v...)
}

// Err log to stderr with timestamp and stack trace
func Err(v ...interface{}) {
	t := time.Now().Format("[2006-01-02 15:04:05]")
	if goos != "windows" {
		v = append(v, "\n"+string(debug.Stack()))
	}
	v = append([]interface{}{C(t, "7")}, v...)

	fmt.Fprintln(Stderr, v...)
}

// ClearScreen ...
func ClearScreen() error {
	if goos == "windows" {
		_, err := os.Stdout.WriteString("\n\n\n\n\n")
		return err
	}

	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

// C color terminal string
func C(str interface{}, color string) string {
	return ansiColor.Color(fmt.Sprint(str), color)
}
