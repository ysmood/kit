package gokit

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/k0kubun/go-ansi"
)

// Stdout ...
var Stdout = ansi.NewAnsiStdout()

// Stderr ...
var Stderr = ansi.NewAnsiStderr()

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

// Dump spew dump
func Dump(v ...interface{}) {
	spew.Dump(v...)
}
