package bootstrap

import "ip-grpc/internal/grpcsrv/server"

func Server() *server.Server {
	return server.NewServer()
}
