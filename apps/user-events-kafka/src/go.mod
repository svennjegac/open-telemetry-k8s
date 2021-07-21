module user-events-kafka

go 1.16

require (
	github.com/confluentinc/confluent-kafka-go v1.7.0
	github.com/svennjegac/opentelemetry-go-contrib/instrumentation/github.com/confluentinc/confluent-kafka-go/kafka/otelkafka v0.0.0-20210716163651-b2f2e00ffe8e
	go.opentelemetry.io/otel v1.0.0-RC1
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.0.0-RC1
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.0.0-RC1
	go.opentelemetry.io/otel/sdk v1.0.0-RC1
	go.opentelemetry.io/otel/trace v1.0.0-RC1
	golang.org/x/net v0.0.0-20210503060351-7fd8e65b6420 // indirect
	golang.org/x/sys v0.0.0-20210616094352-59db8d763f22 // indirect
	google.golang.org/genproto v0.0.0-20210624195500-8bfb893ecb84 // indirect
)
