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
