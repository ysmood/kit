package os_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/ysmood/gokit/pkg/os"
	. "github.com/ysmood/gokit/pkg/utils"
)

func TestOutputString(t *testing.T) {
	str, err := GenerateRandomString(10)
	p := fmt.Sprintf("fixtures/deep/path/%s/output_file", str)

	if err != nil {
		panic(err)
	}

	_ = OutputFile(p, str, nil)

	var c string
	c, err = ReadStringFile(p)

	if err != nil {
		panic(err)
	}

	assert.Equal(t, c, str)
}

func TestOutputBytes(t *testing.T) {
	str, err := GenerateRandomString(10)
	p := fmt.Sprintf("fixtures/deep/path/%s/output_file", str)

	if err != nil {
		panic(err)
	}

	_ = OutputFile(p, []byte("test"), nil)

	var c string
	c, err = ReadStringFile(p)

	if err != nil {
		panic(err)
	}

	assert.Equal(t, c, "test")
}

func TestOutputStringErr(t *testing.T) {
	err := OutputFile("fixtures", "", nil)

	assert.EqualError(t, err, "open fixtures: is a directory")
}

func TestOutputStringErr2(t *testing.T) {
	err := OutputFile("/a/a", "", nil)
	assert.EqualError(t, err, "mkdir /a: permission denied")
}

func TestOutputJSON(t *testing.T) {
	str, err := GenerateRandomString(10)
	p := fmt.Sprintf("fixtures/deep/%s", str)

	if err != nil {
		panic(err)
	}

	data := map[string]interface{}{
		"A": str,
		"B": 10.0,
	}

	_ = OutputFile(p, data, nil)

	var ret interface{}
	err = ReadJSON(p, &ret)

	if err != nil {
		panic(err)
	}

	assert.Equal(t, ret, data)
}

func TestMkdir(t *testing.T) {
	p := "fixtures/deep/a/b/c"
	_ = Mkdir(p, nil)

	assert.Equal(t, true, DirExists(p))
}

func TestGlob(t *testing.T) {
	_ = OutputFile("fixtures/glob/a/b", "", nil)
	_ = OutputFile("fixtures/glob/a/c", "", nil)

	l, err := Walk("glob/**").Dir("fixtures").List()
	E(err)
	assert.Equal(t, 3, len(l))
}

func TestGlobGit(t *testing.T) {
	l, err := Walk("**", WalkGitIgnore).List()
	E(err)
	fullPath, _ := filepath.Abs("fs.go")
	assert.Contains(t, l, fullPath)
}

func TestRemove(t *testing.T) {
	_ = OutputFile("fixtures/remove/a", "", nil)
	_ = OutputFile("fixtures/remove/b/c", "", nil)
	_ = OutputFile("fixtures/remove/b/d", "", nil)
	_ = OutputFile("fixtures/remove/e/f/g", "", nil)

	E(Remove("fixtures/remove/**"))

	l, err := Walk("fixtures/remove/**").List()
	E(err)
	assert.Equal(t, 0, len(l))
}

func TestRemoveSingleFile(t *testing.T) {
	p := "fixtures/remove-single/a"
	_ = OutputFile(p, "", nil)

	assert.Equal(t, true, FileExists(p))

	_ = Remove(p)

	assert.Equal(t, false, FileExists(p))
}

func TestMove(t *testing.T) {
	str, _ := GenerateRandomString(10)
	p := fmt.Sprintf("fixtures/%s", str)

	_ = OutputFile(p+"/a/b", "", nil)
	_ = OutputFile(p+"/a/c", "", nil)

	_ = Move(p+"/a", p+"/d", nil)

	assert.True(t, Exists(p+"/d/b"))
	assert.True(t, DirExists(p+"/d"))
}

func TestGoPath(t *testing.T) {
	s := GoPath()

	assert.True(t, Exists(s))
}
