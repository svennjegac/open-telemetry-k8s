package optelm

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func Setup() func() {
	traceExporter := createOTLPExporter()

	traceProvider, traceProviderShutdown := createTraceProvider(traceExporter)

	setOTELGlobals(traceProvider)

	return traceProviderShutdown
}

func createConsoleExporter() *stdouttrace.Exporter {
	traceExporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
	)
	if err != nil {
		log.Fatalf("failed to initialize stdouttrace export pipeline: %v", err)
	}
	return traceExporter
}

func createOTLPExporter() *otlptrace.Exporter {
	client, err := otlptracegrpc.New(
		context.Background(),
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithTimeout(time.Second*3),
		otlptracegrpc.WithEndpoint(":4317"),
	)
	if err != nil {
		log.Fatalf("failed to initialize grpc trace export pipeline: %v", err)
	}
	return client
}

func createTraceProvider(traceExporter sdktrace.SpanExporter) (*sdktrace.TracerProvider, func()) {
	resources := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("wallet"),
		semconv.ServiceVersionKey.String("2.0"),
		semconv.ServiceInstanceIDKey.String("instance-2xkh78"),
	)

	ctx := context.Background()
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(bsp),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(resources),
	)

	// Handle this error in a sensible manner where possible
	return tp, func() { _ = tp.Shutdown(ctx) }
}

func setOTELGlobals(tp *sdktrace.TracerProvider) {
	otel.SetTracerProvider(tp)
	propagator := propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{})
	otel.SetTextMapPropagator(propagator)
}
