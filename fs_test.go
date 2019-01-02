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

func TestGlob(t *testing.T) {
	g.OutputFile("fixtures/glob/a/b", "", nil)
	g.OutputFile("fixtures/glob/a/c", "", nil)

	l, err := g.Glob([]string{"glob/**"}, &g.WalkOptions{
		Dir: "fixtures",
	})
	g.E(err)
	assert.Equal(t, 3, len(l))
}

func TestRemove(t *testing.T) {
	g.OutputFile("fixtures/remove/a", "", nil)
	g.OutputFile("fixtures/remove/b/c", "", nil)
	g.OutputFile("fixtures/remove/b/d", "", nil)
	g.OutputFile("fixtures/remove/e/f/g", "", nil)

	g.E(g.Remove("fixtures/remove/**"))

	l, err := g.Glob([]string{"fixtures/remove/**"}, nil)
	g.E(err)
	assert.Equal(t, 0, len(l))
}

func TestRemoveSingleFile(t *testing.T) {
	p := "fixtures/remove-single/a"
	g.OutputFile(p, "", nil)

	assert.Equal(t, true, g.FileExists(p))

	g.Remove(p)

	assert.Equal(t, false, g.FileExists(p))
}
