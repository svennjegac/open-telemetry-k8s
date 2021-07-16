package userevents

import (
	"context"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/pkg/errors"
)

const (
	logBrokers = "brokers"
)

type Producer struct {
}

func NewProducer() *Producer {
	return &Producer{}
}

func (p *Producer) Produce(ctx context.Context, id string) error {
	cm := &kafka.ConfigMap{
		"bootstrap.servers":       "192.168.65.2:9092",
		"max.in.flight":           1,
		"socket.keepalive.enable": true,
		"socket.max.fails":        1,

		// "enable.idempotence":           enableIdempotence,
		// "queue.buffering.max.messages": queueBufferingMaxMessages,
		// "queue.buffering.max.kbytes":   queueBufferingMaxKbytes,
		"queue.buffering.max.ms":   5,
		"message.send.max.retries": 2,
		// "retry.backoff.ms":             retryBackoffMs,
		// "compression.codec":            compressionCodec,
		// "batch.num.messages":           batchNumMessages,

		"request.required.acks": -1,
		"request.timeout.ms":    1000,
		"message.timeout.ms":    1000,
		"partitioner":           "murmur2_random", // consistent_random
	}
	producer, err := kafka.NewProducer(cm)
	if err != nil {
		return err
	}

	md, err := producer.GetMetadata(nil, true, 10000)
	if err != nil {
		return err
	} else {
		fmt.Println(md)
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

	delCh := make(chan kafka.Event, 1)
	err = producer.Produce(msg, delCh)
	if err != nil {
		return err
	}

	select {
	case ack := <-delCh:
		fmt.Println(ack)
		return nil
	case <-time.After(time.Second * 3):
		return errors.New("ctx deadline exceeded")
	}

	return err
}
