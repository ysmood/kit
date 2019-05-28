package os_test

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	kit "github.com/ysmood/gokit"
)

func TestMatch(t *testing.T) {
	m, _ := kit.NewMatcher("/root/a", []string{"**", kit.WalkIgnoreHidden})

	matched, negative, _ := m.Match("/root/a/.git", true)
	assert.Equal(t, false, matched)
	assert.Equal(t, true, negative)
}

func TestWalkErr(t *testing.T) {
	oldPath := os.Getenv("PATH")
	os.Setenv("PATH", "")
	defer os.Setenv("PATH", oldPath)

	_, err := kit.Walk("**/file", "!g").List()

	assert.EqualError(t, err, "exec: \"git\": executable file not found in $PATH")
}

func TestWalkGitFatalErr(t *testing.T) {
	_, err := kit.NewMatcher("/", []string{kit.WalkGitIgnore})

	assert.EqualError(t, err, "fatal: not a git repository (or any of the parent directories): .git\nexit status 128")
}

func TestWalkCallbackErr(t *testing.T) {
	err := kit.Walk("**/file", "!g").Do(func(s string, d kit.WalkDirent) error {
		return errors.New("err")
	})

	assert.EqualError(t, err, "err")
}

func TestWalkGitSubmodule(t *testing.T) {
	l := kit.Walk("**/file", "!g").Dir("../../").MustList()

	assert.Len(t, l, 1)
}

func TestWalkParrentGitignore(t *testing.T) {
	p := "tmp/" + kit.GenerateRandomString(10)
	kit.E(kit.OutputFile(p+"/f", "", nil))

	l := kit.Walk(p+"/f", "!g").MustList()

	assert.Len(t, l, 0)
}

func TestWalkOptions(t *testing.T) {
	m, _ := kit.NewMatcher("", []string{"*"})
	l, _ := kit.Walk("*").Sort().FollowSymbolicLinks().Matcher(m).List()

	assert.True(t, len(l) > 0)
}

func TestWalkErrPattern(t *testing.T) {
	assert.EqualError(t, kit.ErrArg(kit.Walk("[]a]").List()), "syntax error in pattern")
}

func TestWalkGitNotFound(t *testing.T) {
	oldPath := os.Getenv("PATH")
	os.Setenv("PATH", "")
	defer os.Setenv("PATH", oldPath)
	_, err := kit.NewMatcher("", []string{"!g"})
	assert.EqualError(t, err, "exec: \"git\": executable file not found in $PATH")
}
