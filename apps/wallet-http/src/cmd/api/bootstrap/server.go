package bootstrap

import (
	"wallet-http/internal/httpsrv/handlers"
	"wallet-http/internal/httpsrv/server"
)

func Server() *server.Server {
	return server.New(handlers.NewWalletRegistrationHandler())
}
