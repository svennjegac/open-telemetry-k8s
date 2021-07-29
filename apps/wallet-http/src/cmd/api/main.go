package main

import (
	"context"

	"wallet-http/cmd/api/bootstrap"
)

func main() {
	// tracerProviderShutdown := optelm.Setup()
	// defer tracerProviderShutdown()

	server := bootstrap.Server()

	server.Start(context.Background())
}
