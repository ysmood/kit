package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ysmood/gokit/pkg/utils"
)

func TestServerErr(t *testing.T) {
	assert.EqualError(t, utils.ErrArg(Server("-1")), "listen tcp: address -1: missing port in address")
}
