package main

import (
	"context"

	"ip-grpc/cmd/api/bootstrap"
	"ip-grpc/internal/optelm"
)

func main() {
	tracerProviderShutdown := optelm.Setup()
	defer tracerProviderShutdown()

	server := bootstrap.Server()

	server.Start(context.Background())
}
