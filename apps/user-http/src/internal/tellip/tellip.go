package tellip

import (
	"context"
	"log"

	"user-http/internal/ip"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
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
	)
	if err != nil {
		log.Fatalln("error dialing target;", err)
	}

	client := ip.NewIPServiceClient(conn)

	return &TellIP{
		tracer: otel.Tracer("sven.njegac/basic"),
		client: client,
	}
}

func (t *TellIP) TellMeYourIP(ctx context.Context) error {
	ctx, span := t.tracer.Start(ctx, "ip-tell-me-ip")
	defer span.End()

	resp, err := t.client.TellMeYourIP(ctx, &ip.TellMeYourIPRequest{
		ClientIp: "user-http-1.1.1.0",
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "network-error")
		return err
	}

	span.AddEvent("got-ip", trace.WithAttributes(attribute.String("ip", resp.ServerIp)))

	return nil
}
