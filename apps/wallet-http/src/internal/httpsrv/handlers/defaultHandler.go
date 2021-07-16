package handlers

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type DefaultHandler struct {
	tracer      trace.Tracer
	propagators propagation.TextMapPropagator
}

func NewDefaultHandler() *DefaultHandler {
	return &DefaultHandler{
		tracer:      otel.Tracer("sven.njegac/basic-2"),
		propagators: otel.GetTextMapPropagator(),
	}
}

func (d *DefaultHandler) Default() httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

		fmt.Println("###### Handling new request")

		fmt.Println("Top level context")
		fmt.Printf("%s\n", request.Context())
		fmt.Printf("%+v\n", request.Context())
		fmt.Printf("%p\n", request.Context())
		fmt.Printf("%T\n", request.Context())

		var hf http.HandlerFunc

		hf = func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Inside level context")
			fmt.Printf("%s\n", r.Context())
			fmt.Printf("%+v\n", r.Context())
			fmt.Printf("%p\n", r.Context())
			fmt.Printf("%T\n", r.Context())

			ctx := d.propagators.Extract(r.Context(), propagation.HeaderCarrier(r.Header))

			ctx, span := d.tracer.Start(ctx, "register-wallet")
			defer span.End()

			time.Sleep(time.Duration(rand.Intn(50)+30) * time.Millisecond)

			// bag := baggage.FromContext(ctx)
			// fmt.Println(bag.String())
			// fmt.Println(bag.Len())
			// fmt.Println(bag.Members())

			w.WriteHeader(http.StatusOK)
			_, err := fmt.Fprintf(w, "okeish")
			if err != nil {
				log.Println("default err", err)
			}
		}

		h := otelhttp.NewHandler(hf, "middleware-otel-register-wallet")
		h.ServeHTTP(writer, request)
	}
}

func (d *DefaultHandler) Hello() httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {

		writer.WriteHeader(http.StatusOK)
		_, err := fmt.Fprintf(writer, "hello route, time: %s", time.Now())
		if err != nil {
			log.Println("default err", err)
		}
	}
}
