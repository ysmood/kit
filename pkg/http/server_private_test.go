package http

import (
	"testing"
)

func TestServerErr(t *testing.T) {
	_, err := Server("-1")
	if err.Error() != "listen tcp: address -1: missing port in address" {
		panic(err)
	}
}
