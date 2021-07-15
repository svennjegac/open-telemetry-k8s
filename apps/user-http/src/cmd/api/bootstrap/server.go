package bootstrap

import (
	"user-http/internal/httpsrv/handlers"
	"user-http/internal/httpsrv/server"
)

func Server() *server.Server {
	return server.New(handlers.NewDefaultHandler())
}
