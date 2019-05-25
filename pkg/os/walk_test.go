package os_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/ysmood/gokit/pkg/os"
	. "github.com/ysmood/gokit/pkg/utils"
)

func TestMatch(t *testing.T) {
	m, _ := NewMatcher("/root/a", []string{"**", WalkIgnoreHidden})

	matched, negative, _ := m.Match("/root/a/.git", true)
	assert.Equal(t, false, matched)
	assert.Equal(t, true, negative)
}

func TestWalkGitSubmodule(t *testing.T) {
	l := Walk("**/file", "!g").Dir("../../").MustList()

	assert.Len(t, l, 1)
}

func TestWalkParrentGitignore(t *testing.T) {
	p := "tmp/" + GenerateRandomString(10)
	E(OutputFile(p+"/f", "", nil))

	l := Walk(p+"/f", "!g").MustList()

	assert.Len(t, l, 0)
}

func TestWalkOptions(t *testing.T) {
	m, _ := NewMatcher("", []string{"*"})
	l, _ := Walk("*").Sort().FollowSymbolicLinks().Matcher(m).List()

	assert.True(t, len(l) > 0)
}

func TestWalkErrPattern(t *testing.T) {
	assert.EqualError(t, ErrArg(Walk("[]a]").List()), "syntax error in pattern")
}
