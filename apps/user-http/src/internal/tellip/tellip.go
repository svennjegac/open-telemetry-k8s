package tellip

import (
	"context"
	"log"
	"time"

	"user-http/internal/ip"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type TellIP struct {
	tracer trace.Tracer
	client ip.IPServiceClient
}

func NewTellIP() *TellIP {
	conn, err := grpc.DialContext(
		context.Background(),
		"dns:///ip-grpc.default.svc.cluster.local:8113",
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithDefaultServiceConfig("{\"loadBalancingPolicy\":\"round_robin\"}"),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		log.Fatalln("error dialing target;", err)
	}

	client := ip.NewIPServiceClient(conn)

	return &TellIP{
		tracer: otel.Tracer("sven.njegac/open-telemetry-k8s"),
		client: client,
	}
}

func (t *TellIP) TellMeYourIP(ctx context.Context) error {
	ctx, span := t.tracer.Start(ctx, "tell-me-your-ip")
	defer span.End()

	md := metadata.Pairs(
		"timestamp", time.Now().Format(time.StampNano),
		"client-id", "web-api-client-us-east-1",
		"user-id", "some-test-user-id",
	)

	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := t.client.TellMeYourIP(ctx, &ip.TellMeYourIPRequest{
		ClientIp: "user-http-1.0.40",
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "network-error")
		return err
	}

	span.AddEvent("got-ip", trace.WithAttributes(attribute.String("ip", resp.ServerIp)))

	return nil
}
