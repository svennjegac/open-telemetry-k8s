package wallet

import (
	"context"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Wallet struct {
	tracer trace.Tracer
	client *http.Client
}

func New() *Wallet {
	cli := &http.Client{}
	cli.Transport = otelhttp.NewTransport(cli.Transport)
	return &Wallet{
		tracer: otel.Tracer("sven.njegac/open-telemetry-k8s"),
		client: cli,
	}
}

func (h *Wallet) RegisterToWallet(ctx context.Context, id string) error {
	ctx, span := h.tracer.Start(ctx, "register-wallet")
	defer span.End()

	prop, err := baggage.NewKeyValueProperty("id-property", id)
	if err != nil {
		span.SetStatus(codes.Error, "property construction")
		span.RecordError(err)
		return err
	}

	idMember, err := baggage.NewMember("id", id, prop)
	if err != nil {
		span.SetStatus(codes.Error, "member construction")
		span.RecordError(err)
		return err
	}

	bag, err := baggage.New(idMember)
	if err != nil {
		span.SetStatus(codes.Error, "bag construction")
		span.RecordError(err)
		return err
	}

	ctx = baggage.ContextWithBaggage(ctx, bag)

	// make a new http request
	r, err := http.NewRequest("GET", "http://wallet-http.default.svc.cluster.local:8112/register", nil)
	if err != nil {
		panic(err)
	}
	r = r.WithContext(ctx)

	resp, err := h.client.Do(r)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "http-call")
		return err
	}
	resp.Body.Close()

	return nil
}
