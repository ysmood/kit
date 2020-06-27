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

	server *http.Server
}

// GinContext ...
type GinContext = *gin.Context

// Server listen to a port then create a gin server.
// I created this wrapper because gin doesn't give a signal to tell when the
// port is ready.
func Server(address string) (*ServerContext, error) {
	s := &ServerContext{
		server: &http.Server{},
	}

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

// Set options
func (ctx *ServerContext) Set(server *http.Server) *ServerContext {
	ctx.server = server
	return ctx
}

// Do start the handler loop
func (ctx *ServerContext) Do() error {
	ctx.server.Handler = ctx.Engine
	return ctx.server.Serve(ctx.Listener)
}

// MustDo ...
func (ctx *ServerContext) MustDo() {
	utils.E(ctx.Do())
}
