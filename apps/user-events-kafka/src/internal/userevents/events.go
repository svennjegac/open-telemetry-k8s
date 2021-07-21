package userevents

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/svennjegac/opentelemetry-go-contrib/instrumentation/github.com/confluentinc/confluent-kafka-go/kafka/otelkafka"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Consumer struct {
	tracer trace.Tracer
}

func NewConsumer() *Consumer {
	return &Consumer{
		tracer: otel.Tracer("sven.njegac/open-telemetry-k8s"),
	}
}

func (c *Consumer) Consume(ctx context.Context) error {
	cm := &kafka.ConfigMap{
		"queue.buffering.max.ms":   5,
		"message.send.max.retries": 2,
		"request.required.acks":    -1,
		"request.timeout.ms":       1000,
		"message.timeout.ms":       1000,
		"partitioner":              "murmur2_random", // consistent_random
		"bootstrap.servers":        "my-cluster-kafka-brokers.kafka.svc.cluster.local:9092",
		"socket.keepalive.enable":  true,
		"socket.max.fails":         1,
		"group.id":                 "user-events-otel-consumer",
		"session.timeout.ms":       5000,
		"heartbeat.interval.ms":    3000,
		"max.poll.interval.ms":     6000,
		"enable.auto.commit":       true,
		"auto.commit.interval.ms":  5000,
		"go.events.channel.enable": false,
		"enable.auto.offset.store": false,
	}

	consumer, err := kafka.NewConsumer(cm)
	if err != nil {
		return err
	}

	_, err = consumer.GetMetadata(nil, true, 10000)
	if err != nil {
		return err
	}

	topic := "user-events-otel"
	consumer.Assign([]kafka.TopicPartition{{
		Topic:     &topic,
		Partition: 0,
		Offset:    kafka.OffsetStored,
	}})

	cs := otelkafka.WrapConsumer(consumer)

	for {
		event := cs.Poll(3000)
		if event == nil {
			continue
		}

		switch e := event.(type) {
		case *kafka.Message:
			fmt.Println("Got new message")
			fmt.Println("Key", string(e.Key))
			fmt.Println("Value", string(e.Value))
			fmt.Println("---")

			ctx = otel.GetTextMapPropagator().Extract(ctx, otelkafka.NewMessageCarrier(e))

			ctx, span := c.tracer.Start(ctx, "consumed-message")
			time.Sleep(time.Duration(rand.Intn(30)+30) * time.Millisecond)

			ctx, span2 := c.tracer.Start(ctx, "writing-message-to-mongo")
			time.Sleep(time.Duration(rand.Intn(30)+30) * time.Millisecond)
			span2.End()

			span.End()

		default:
			fmt.Println(event)
		}
	}
}
