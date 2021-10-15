package wallet

import (
	"context"
	base642 "encoding/base64"
	"log"
	"net/http"
	"time"

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
	// Start new span from parent context.
	ctx, span := h.tracer.Start(ctx, "register-wallet")
	defer span.End()

	// Example: Add some baggage.
	baggageVal, err := h.reqBaggageVal(id)
	if err != nil {
		// This sets:
		// - span status code to error
		// - span status message to "baggage construction"
		span.SetStatus(codes.Error, "baggage construction")
		// This adds log to span (event), but does not set span status to error.
		span.RecordError(err)
		return err
	}

	ctx = baggage.ContextWithBaggage(ctx, baggageVal)

	// make a new http request
	r, err := http.NewRequest("GET", "http://wallet-http.default.svc.cluster.local:8112/register-user", nil)
	if err != nil {
		return err
	}
	r = r.WithContext(ctx)

	resp, err := h.client.Do(r)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "http call")
		return err
	}
	resp.Body.Close()

	return nil
}

func (h *Wallet) reqBaggageVal(id string) (baggage.Baggage, error) {
	idProp, err := baggage.NewKeyValueProperty("id-property", id)
	if err != nil {
		log.Println("id property", err)
		return baggage.Baggage{}, err
	}

	idHashProp, err := baggage.NewKeyValueProperty("id-property-base64", base64(id))
	if err != nil {
		log.Println("id property base64", err)
		return baggage.Baggage{}, err
	}

	idMember, err := baggage.NewMember("id-baggage-key", "id-baggage-value", idProp, idHashProp)
	if err != nil {
		log.Println("id member", err)
		return baggage.Baggage{}, err
	}

	timeProp, err := baggage.NewKeyValueProperty("time-property", base64(time.Now().String()))
	if err != nil {
		log.Println("time property", err)
		return baggage.Baggage{}, err
	}

	timeMember, err := baggage.NewMember("time-baggage-key", "time-baggage-value", timeProp)
	if err != nil {
		log.Println("time member", err)
		return baggage.Baggage{}, err
	}

	bag, err := baggage.New(idMember, timeMember)
	if err != nil {
		log.Println("new baggage", err)
		return baggage.Baggage{}, err
	}

	return bag, nil
}

func base64(s string) string {
	return base642.StdEncoding.EncodeToString([]byte(s))
}
