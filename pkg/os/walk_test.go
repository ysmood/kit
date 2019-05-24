package os_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/ysmood/gokit/pkg/os"
)

func TestMatch(t *testing.T) {
	m, _ := NewMatcher("/root/a", []string{"**", WalkIgnoreHidden})

	matched, negative, _ := m.Match("/root/a/.git", true)
	assert.Equal(t, false, matched)
	assert.Equal(t, true, negative)
}

func TestWalkGitSubmodule(t *testing.T) {
	l, _ := Walk("**/file", "!g").Dir("../../").List()

	assert.Len(t, l, 1)
}

func TestWalkOptions(t *testing.T) {
	m, _ := NewMatcher("", []string{"*"})
	l, _ := Walk("*").Sort().FollowSymbolicLinks().Matcher(m).List()

	assert.True(t, len(l) > 0)
}
