package gokit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExec(t *testing.T) {
	err := Exec("go", "version")
	assert.Equal(t, nil, err)
}
func TestExecPrefix(t *testing.T) {
	err := Exec("echo", "test", ExecOptions{
		Prefix: "[app] ",
	})
	assert.Equal(t, nil, err)
}

func TestExecErr(t *testing.T) {
	err := Exec("exitexit")
	assert.NotEqual(t, nil, err)
}
