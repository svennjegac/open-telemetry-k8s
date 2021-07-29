package main

import (
	"context"
	"fmt"

	"user-events-kafka/cmd/api/bootstrap"
)

func main() {
	// tracerProviderShutdown := optelm.Setup()
	// defer tracerProviderShutdown()

	consumer := bootstrap.Consumer()

	fmt.Println(consumer.Consume(context.Background()))
}
