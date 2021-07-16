module user-http

go 1.16

require (
	cloud.google.com/go/logging v1.4.2
	github.com/cenkalti/backoff/v4 v4.1.1
	github.com/confluentinc/confluent-kafka-go v1.7.0 // indirect
	github.com/julienschmidt/httprouter v1.3.0
	github.com/pkg/errors v0.9.1
	github.com/svennjegac/opentelemetry-go-contrib/instrumentation/github.com/confluentinc/confluent-kafka-go/kafka/otelkafka v0.0.0-20210716142334-4887023060e0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.21.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.21.0
	go.opentelemetry.io/otel v1.0.0-RC1
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.0.0-RC1
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.0.0-RC1
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.0.0-RC1
	go.opentelemetry.io/otel/sdk v1.0.0-RC1
	go.opentelemetry.io/otel/trace v1.0.0-RC1
	google.golang.org/api v0.50.0 // indirect
	google.golang.org/grpc v1.38.0
	google.golang.org/protobuf v1.26.0
)
