package os_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/ysmood/gokit/pkg/os"
	. "github.com/ysmood/gokit/pkg/utils"
)

func TestOutputString(t *testing.T) {
	p := "tmp/" + GenerateRandomString(10)

	_ = OutputFile(p, p, nil)

	c, err := ReadStringFile(p)

	if err != nil {
		panic(err)
	}

	assert.Equal(t, c, p)
}

func TestOutputBytes(t *testing.T) {
	p := "tmp/" + GenerateRandomString(10)

	_ = OutputFile(p, []byte("test"), nil)

	c, err := ReadStringFile(p)

	if err != nil {
		panic(err)
	}

	assert.Equal(t, c, "test")
}

func TestOutputStringErr(t *testing.T) {
	err := OutputFile(".", "", nil)

	assert.EqualError(t, err, "open .: is a directory")
}

func TestOutputStringErr2(t *testing.T) {
	err := OutputFile("/a/a", "", nil)
	assert.EqualError(t, err, "mkdir /a: permission denied")
}

func TestOutputJSON(t *testing.T) {
	p := "tmp/deep/" + GenerateRandomString(10)

	data := map[string]interface{}{
		"A": p,
		"B": 10.0,
	}

	_ = OutputFile(p, data, nil)

	var ret interface{}
	err := ReadJSON(p, &ret)

	if err != nil {
		panic(err)
	}

	assert.Equal(t, ret, data)
}

func TestMkdir(t *testing.T) {
	p := "tmp/deep/a/b/c"
	_ = Mkdir(p, nil)

	assert.Equal(t, true, DirExists(p))
}

func TestGlob(t *testing.T) {
	_ = OutputFile("tmp/glob/a/b", "", nil)
	_ = OutputFile("tmp/glob/a/c", "", nil)

	l, err := Walk("glob/**").Dir("tmp").List()
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
	_ = OutputFile("tmp/remove/a", "", nil)
	_ = OutputFile("tmp/remove/b/c", "", nil)
	_ = OutputFile("tmp/remove/b/d", "", nil)
	_ = OutputFile("tmp/remove/e/f/g", "", nil)

	E(Remove("tmp/remove/**"))

	l, err := Walk("tmp/remove/**").List()
	E(err)
	assert.Equal(t, 0, len(l))
}

func TestRemoveSingleFile(t *testing.T) {
	p := "tmp/remove-single/a"
	_ = OutputFile(p, "", nil)

	assert.Equal(t, true, FileExists(p))

	_ = Remove(p)

	assert.Equal(t, false, FileExists(p))
}

func TestMove(t *testing.T) {
	p := "tmp/" + GenerateRandomString(10)

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
