package gokit_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	g "github.com/ysmood/gokit"
)

func TestAll(t *testing.T) {
	g.All(func() {
		fmt.Println("one")
	}, func() {
		fmt.Println("two")
	})
}

func TestE(t *testing.T) {
	defer func() {
		r := recover()

		assert.Equal(t, "exec: \"exitexit\": executable file not found in $PATH", r.(error).Error())
	}()

	g.E(g.Exec([]string{"exitexit"}, nil))
}
