package gokit_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	g "github.com/ysmood/gokit"
)

func TestExec(t *testing.T) {
	_, err := g.Exec([]string{"go", "version"}, nil)
	assert.Equal(t, nil, err)
}
func TestExecPrefix(t *testing.T) {
	_, err := g.Exec([]string{"echo", "test"}, &g.ExecOptions{
		Prefix: "[app] ",
	})
	assert.Equal(t, nil, err)
}

func TestExecErr(t *testing.T) {
	_, err := g.Exec([]string{"exitexit"}, nil)
	assert.NotEqual(t, nil, err)
}
