package userevents

import (
	"context"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/pkg/errors"
	"github.com/svennjegac/opentelemetry-go-contrib/instrumentation/github.com/confluentinc/confluent-kafka-go/kafka/otelkafka"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

const (
	logBrokers = "brokers"
)

type Producer struct {
	tracer trace.Tracer
}

func NewProducer() *Producer {
	return &Producer{
		tracer: otel.Tracer("sven.njegac/basic"),
	}
}

func (p *Producer) Produce(ctx context.Context, id string) error {
	ctx, span := p.tracer.Start(ctx, "user-events-in-app-produce")
	defer span.End()

	cm := &kafka.ConfigMap{
		"bootstrap.servers":       "my-cluster-kafka-brokers.kafka.svc.cluster.local:9092",
		"max.in.flight":           1,
		"socket.keepalive.enable": true,
		"socket.max.fails":        1,

		// "enable.idempotence":           enableIdempotence,
		// "queue.buffering.max.messages": queueBufferingMaxMessages,
		// "queue.buffering.max.kbytes":   queueBufferingMaxKbytes,
		"queue.buffering.max.ms":   5,
		"message.send.max.retries": 3,
		// "retry.backoff.ms":             retryBackoffMs,
		// "compression.codec":            compressionCodec,
		// "batch.num.messages":           batchNumMessages,

		"request.required.acks": -1,
		"request.timeout.ms":    4000,
		"message.timeout.ms":    4000,
		"partitioner":           "murmur2_random", // consistent_random
	}
	producer, err := kafka.NewProducer(cm)
	if err != nil {
		span.RecordError(err)
		return errors.Wrap(err, "new producer")
	}


	_, err = producer.GetMetadata(nil, true, 10000)
	if err != nil {
		span.RecordError(err)
		return errors.Wrap(err, "get metadata")
	}

	// Drain events channel to prevent memory leak
	go func() {
		for event := range producer.Events() {
			if _, ok := event.(*kafka.Stats); ok {

			}
		}
	}()

	topic := "user-events-otel"
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: 0,
			Offset:    0,
			Metadata:  nil,
			Error:     nil,
		},
		Value:         []byte(id),
		Key:           []byte(id),
		Timestamp:     time.Time{},
		TimestampType: 0,
		Opaque:        nil,
		Headers:       nil,
	}

	pr := otelkafka.WrapProducer(producer)

	carrier := otelkafka.NewMessageCarrier(msg)

	otel.GetTextMapPropagator().Inject(ctx, carrier)

	delCh := make(chan kafka.Event, 1)
	err = pr.Produce(msg, delCh)
	if err != nil {
		span.RecordError(err)
		return errors.Wrap(err, "produce")
	}

	select {
	case ack := <-delCh:
		fmt.Println(ack)
		return nil
	case <-time.After(time.Second * 5):
		span.RecordError(errors.New("ctx deadline exceeded"))
		return errors.New("ctx deadline exceeded")
	}

	return err
}
