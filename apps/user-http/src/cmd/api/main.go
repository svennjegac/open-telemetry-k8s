package main

import (
	"context"

	"user-http/cmd/api/bootstrap"
	"user-http/internal/optelm"
)

func main() {
	optelm.Setup()

	server := bootstrap.Server()

	server.Start(context.Background())
}
