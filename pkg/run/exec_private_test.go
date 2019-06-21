package run

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testReader struct{}

func (r *testReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("err")
}

func TestPipeToStdoutWithPrefixErr(t *testing.T) {
	r := &testReader{}
	assert.EqualError(t, pipeToStdoutWithPrefix("prefix", r), "err")
}
