package main

import (
	"fmt"

	"user-events-kafka/internal/optelm"
	"user-events-kafka/internal/userevents"
)

func main() {
	optelm.Setup()

	consumer := userevents.NewConsumer()

	err := consumer.Consume()
	fmt.Println(err)
}
