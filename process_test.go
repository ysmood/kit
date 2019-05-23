package gokit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExec(t *testing.T) {
	err := Exec("go", "version").Do()
	assert.Nil(t, err)
}

func TestExecPrefix(t *testing.T) {
	err := Exec("echo", "test").Prefix("[app] ").Do()
	assert.Nil(t, err)
}

func TestExecPrefixColor(t *testing.T) {
	err := Exec("echo", "test").Prefix("[app] @green").Do()
	assert.Nil(t, err)
}

func TestExecErr(t *testing.T) {
	err := Exec("exitexit")
	assert.EqualError(t, err.Do(), "exec: \"exitexit\": executable file not found in $PATH")
}

func TestExecRaw(t *testing.T) {
	err := Exec("echo", "ok").Raw().Do()
	assert.Nil(t, err)
}
