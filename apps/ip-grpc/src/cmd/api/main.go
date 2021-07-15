package main

import (
	"context"

	"ip-grpc/cmd/api/bootstrap"
	"ip-grpc/internal/optelm"
)

func main() {
	optelm.Setup()

	server := bootstrap.Server()

	server.Start(context.Background())
}
