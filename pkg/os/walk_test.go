package os_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	kit "github.com/ysmood/gokit"
)

func TestMatch(t *testing.T) {
	gitIgnorePath := kit.HomeDir() + "/.gitignore_global"
	if !kit.FileExists(gitIgnorePath) {
		kit.E(kit.OutputFile(gitIgnorePath, "", nil))
	}

	m := kit.NewMatcher("/root/a", []string{"**", kit.WalkIgnoreHidden})

	matched, negative, _ := m.Match("/root/a/.git", true)
	assert.Equal(t, false, matched)
	assert.Equal(t, true, negative)
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
	p := "tmp/" + kit.RandString(10)
	kit.E(kit.OutputFile(p+"/f", "", nil))

	l := kit.Walk(p+"/f", "!g").MustList()

	assert.Len(t, l, 0)
}

func TestWalkOptions(t *testing.T) {
	m := kit.NewMatcher("", []string{"*"})
	l, _ := kit.Walk("*").Sort().FollowSymbolicLinks().Matcher(m).List()

	assert.True(t, len(l) > 0)
}

func TestWalkErrPattern(t *testing.T) {
	assert.EqualError(t, kit.ErrArg(kit.Walk("[]a]").List()), "syntax error in pattern")
}
