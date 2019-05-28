package run

import (
	"testing"

	"golang.org/x/crypto/ssh/terminal"
)

func TestRestoreState(t *testing.T) {
	restoreState(&terminal.State{})
}
