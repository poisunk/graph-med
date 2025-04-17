package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"graph-med/internal/base/conf"
	"graph-med/internal/middleware"
	"graph-med/internal/router"
)

type Server struct {
	engine *gin.Engine
	router *router.Router
	addr   string
}

func NewServer(config *conf.AllConfig, router *router.Router) *Server {
	serverConf := &config.Server

	engine := gin.Default()
	engine.Use(middleware.Cors())
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	addr := fmt.Sprintf(":%d", serverConf.Port)

	server := &Server{
		engine: engine,
		router: router,
		addr:   addr,
	}

	server.registerRoutes()
	return server
}

func (s *Server) Run() error {
	return s.engine.Run(s.addr)
}

func (s *Server) registerRoutes() {
	s.router.RegisterRoutes(s.engine)
}
