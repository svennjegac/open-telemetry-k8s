package server

import (
	"github.com/julienschmidt/httprouter"
)

func (s *Server) setRoutes() {
	router := httprouter.New()

	router.GET("/users/:id", s.userHandler.GetUser())

	s.httpServer.Handler = router
}
