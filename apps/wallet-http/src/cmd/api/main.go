package main

import (
	"context"

	"wallet-http/cmd/api/bootstrap"
	"wallet-http/internal/optelm"
)

func main() {
	tracerProviderShutdown := optelm.Setup()
	defer tracerProviderShutdown()

	server := bootstrap.Server()

	server.Start(context.Background())
}
