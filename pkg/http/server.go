package http

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ysmood/gokit/pkg/utils"
)

type ServerContext struct {
	Handler  *gin.Engine
	Listener net.Listener
}

type GinContext = *gin.Context

func Server(address string) (*ServerContext, error) {
	s := &ServerContext{}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	ln, err := net.Listen("tcp", address)

	if err != nil {
		return nil, err
	}

	s.Handler = r
	s.Listener = ln

	return s, nil
}

func MustServer(address string) *ServerContext {
	return utils.E(Server(address))[0].(*ServerContext)
}

func (ctx *ServerContext) Do() error {
	return http.Serve(ctx.Listener, ctx.Handler)
}

func (ctx *ServerContext) MustDo() {
	utils.E(ctx.Do())
}
