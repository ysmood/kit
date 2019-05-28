package http

import (
	"errors"
	"io"
	"testing"
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

func TestReadBodyReadErr(t *testing.T) {
	obj := ErrReader{
		readErr: errors.New("err"),
	}

	_, err := readBody(obj)

	if err.Error() != "err" {
		panic(err)
	}
}

func TestReadBodyCloseErr(t *testing.T) {
	obj := ErrReader{
		closeErr: errors.New("err"),
	}

	_, err := readBody(obj)

	if err.Error() != "err" {
		panic(err)
	}
}
