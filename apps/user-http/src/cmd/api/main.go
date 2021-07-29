package main

import (
	"context"

	"user-http/cmd/api/bootstrap"
)

func main() {
	// tracerProviderShutdown := optelm.Setup()
	// defer tracerProviderShutdown()

	server := bootstrap.Server()

	server.Start(context.Background())
}
