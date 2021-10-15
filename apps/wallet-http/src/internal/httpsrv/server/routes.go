package server

import (
	"github.com/julienschmidt/httprouter"
)

func (s *Server) setRoutes() {
	router := httprouter.New()

	router.GET("/register-user", s.walletRegistrationHandler.RegisterUser())

	s.httpServer.Handler = router
}
