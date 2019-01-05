package gokit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatch(t *testing.T) {
	m, _ := newMatcher("/root/a", []string{"**", WalkHidden})

	matched, negative, _ := m.match("/root/a/.git", true)
	assert.Equal(t, false, matched)
	assert.Equal(t, true, negative)
}
