package http

import (
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ysmood/kit/pkg/utils"
)

// ServerContext ...
type ServerContext struct {
	Engine   *gin.Engine
	Listener net.Listener
}

// GinContext ...
type GinContext = *gin.Context

// Server listen to a port then create a gin server.
// I created this wrapper because gin doesn't give a signal to tell when the
// port is ready.
func Server(address string) (*ServerContext, error) {
	s := &ServerContext{}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	ln, err := net.Listen("tcp", address)

	if err != nil {
		return nil, err
	}

	s.Engine = r
	s.Listener = ln

	return s, nil
}

// MustServer ...
func MustServer(address string) *ServerContext {
	return utils.E(Server(address))[0].(*ServerContext)
}

// Do start the handler loop
func (ctx *ServerContext) Do() error {
	return http.Serve(ctx.Listener, ctx.Engine)
}

// MustDo ...
func (ctx *ServerContext) MustDo() {
	utils.E(ctx.Do())
}
