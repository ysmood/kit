package gokit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExec(t *testing.T) {
	err := Exec([]string{"go", "version"}, nil)
	assert.Equal(t, nil, err)
}
func TestExecPrefix(t *testing.T) {
	err := Exec([]string{"echo", "test"}, &ExecOptions{
		Prefix: "[app] ",
	})
	assert.Equal(t, nil, err)
}

func TestExecErr(t *testing.T) {
	err := Exec([]string{"exitexit"}, nil)
	assert.NotEqual(t, nil, err)
}
