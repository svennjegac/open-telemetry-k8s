package handlers

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/trace"
)

type WalletRegistrationHandler struct {
	tracer trace.Tracer
}

func NewWalletRegistrationHandler() *WalletRegistrationHandler {
	return &WalletRegistrationHandler{
		tracer: otel.Tracer("sven.njegac/open-telemetry-k8s"),
	}
}

func (d *WalletRegistrationHandler) RegisterUser() httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

		var hf http.HandlerFunc
		hf = func(w http.ResponseWriter, r *http.Request) {
			log.Println("- - - - - - - - - - - - - - -")
			log.Println("new request processing started")

			// Start new span from parent context.
			ctx, span := d.tracer.Start(r.Context(), "register-user-handler")
			defer span.End()

			// Fake delay.
			time.Sleep(time.Duration(rand.Intn(50)+30) * time.Millisecond)

			// Example - logging of received baggage.
			d.logBaggage(ctx)

			d.logRequestHeaders(request)

			w.WriteHeader(http.StatusOK)
			_, err := fmt.Fprintf(w, "okay")
			if err != nil {
				log.Println("response writer err", err)
			}
		}

		h := otelhttp.NewHandler(hf, "register-user-middleware")
		h.ServeHTTP(writer, request)
	}
}

func (d *WalletRegistrationHandler) logRequestHeaders(request *http.Request) {
	log.Println("$ $ $ $ $ $ $ $ $ $ $ $ $ $ $ $ $ $ $ $ $ $ $ $ $ $ $ $ $ $")
	log.Println("printing headers")
	for key, value := range request.Header {
		log.Println("new header")
		log.Println("key:", key)
		log.Println("value:", value)
	}
}

func (d *WalletRegistrationHandler) logBaggage(ctx context.Context) {
	// Extract and print baggage info.
	// This is just an example how to use baggage.
	bag := baggage.FromContext(ctx)
	log.Println("# # # # # # # # # # # # # # # # # # # # # # # # # # # # # #")
	log.Println("printing baggage")
	for _, mem := range bag.Members() {
		log.Println("new member...")
		for _, prop := range mem.Properties() {
			log.Println("propKey", prop.Key())
			val, ok := prop.Value()
			log.Println("propValue", val, ok)
		}
	}
}
