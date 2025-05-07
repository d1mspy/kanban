package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	host string
}

func New(host string) *Server {
	s := &Server{host: host}
	return s
}

func (r *Server) newAPI() *gin.Engine {
	engine := gin.New()

	engine.GET("/health", func(ctx *gin.Context) {
		ctx.Status(http.StatusNoContent)
	})

	return engine
}

func (r *Server) Start() {
	r.newAPI().Run(r.host)
}