package gokit_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	g "github.com/ysmood/gokit"
)

func TestOutputString(t *testing.T) {
	str, err := g.GenerateRandomString(10)
	p := fmt.Sprintf("fixtures/deep/path/%s/output_file", str)

	if err != nil {
		panic(err)
	}

	g.OutputFile(p, str, nil)

	var c string
	c, err = g.ReadStringFile(p)

	if err != nil {
		panic(err)
	}

	assert.Equal(t, c, str)
}
func TestOutputJSON(t *testing.T) {
	str, err := g.GenerateRandomString(10)
	p := fmt.Sprintf("fixtures/deep/%s", str)

	if err != nil {
		panic(err)
	}

	data := map[string]interface{}{
		"A": str,
		"B": 10.0,
	}

	g.OutputFile(p, data, nil)

	var ret interface{}
	err = g.ReadJSON(p, &ret)

	if err != nil {
		panic(err)
	}

	assert.Equal(t, ret, data)
}

func TestMkdir(t *testing.T) {
	p := "fixtures/deep/a/b/c"
	g.Mkdir(p, nil)

	assert.Equal(t, true, g.DirExists(p))
}