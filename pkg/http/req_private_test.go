package http

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ErrReader struct {
	readErr  error
	closeErr error
}

func (r ErrReader) Read(p []byte) (int, error) {
	if r.readErr != nil {
		return 0, r.readErr
	}

	return 0, io.EOF
}

func (r ErrReader) Close() error {
	return r.closeErr
}

func TestReadBodyReadErro(t *testing.T) {
	obj := ErrReader{
		readErr: errors.New("err"),
	}

	_, err := readBody(obj)

	assert.EqualError(t, err, "err")
}

func TestReadBodyCloseErro(t *testing.T) {
	obj := ErrReader{
		closeErr: errors.New("err"),
	}

	_, err := readBody(obj)

	assert.EqualError(t, err, "err")
}
