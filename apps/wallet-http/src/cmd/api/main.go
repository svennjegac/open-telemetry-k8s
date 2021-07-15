package main

import (
	"context"

	"wallet-http/cmd/api/bootstrap"
	"wallet-http/internal/optelm"
)

func main() {
	optelm.Setup()

	server := bootstrap.Server()

	server.Start(context.Background())
}
