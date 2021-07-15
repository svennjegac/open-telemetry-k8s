package server

import (
	"github.com/julienschmidt/httprouter"
)

func (s *Server) setRoutes() {
	router := httprouter.New()

	router.GET("/users/:id", s.defaultHandler.Default())
	router.GET("/hello", s.defaultHandler.Hello())

	s.httpServer.Handler = router
}
