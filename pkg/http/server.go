package http

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ServerContext struct {
	Handler  *gin.Engine
	Listener net.Listener
	Error    error
}

type GinContext = *gin.Context

func Server(address string) *ServerContext {
	s := &ServerContext{}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	ln, err := net.Listen("tcp", address)

	if err != nil {
		s.Error = err
		return s
	}

	s.Handler = r
	s.Listener = ln

	return s
}

func (ctx *ServerContext) Do() {
	ctx.Error = http.Serve(ctx.Listener, ctx.Handler)
}
