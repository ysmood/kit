package gokit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExec(t *testing.T) {
	err := Exec("go", "version").Do()
	assert.Equal(t, nil, err)
}

func TestExecPrefix(t *testing.T) {
	err := Exec("echo", "test").Prefix("[app] ").Do()
	assert.Equal(t, nil, err)
}

func TestExecPrefixColor(t *testing.T) {
	err := Exec("echo", "test").Prefix("[app] @green").Do()
	assert.Equal(t, nil, err)
}

func TestExecErr(t *testing.T) {
	err := Exec("exitexit")
	assert.NotEqual(t, nil, err)
}
