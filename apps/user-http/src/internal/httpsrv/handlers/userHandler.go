package handlers

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"user-http/internal/memorydb"
	"user-http/internal/models"
	"user-http/internal/tellip"
	"user-http/internal/userevents"
	"user-http/internal/wallet"

	"github.com/julienschmidt/httprouter"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type UserHandler struct {
	tracer                 trace.Tracer
	otelUserIDKey          attribute.Key
	otelTimeKey            attribute.Key
	userRepository         UserRepository
	walletRepository       WalletRepository
	tellMeYourIPRepository TellMeYourIPRepository
	userEventsProducer     UserEventsProducer
}

func NewDefaultHandler() *UserHandler {
	return &UserHandler{
		tracer:                 otel.Tracer("sven.njegac/open-telemetry-k8s"),
		otelUserIDKey:          "get-user/id",
		otelTimeKey:            "get-user/time",
		userRepository:         memorydb.New(),
		walletRepository:       wallet.New(),
		tellMeYourIPRepository: tellip.NewTellIP(),
		userEventsProducer:     userevents.NewProducer(),
	}
}

func (d *UserHandler) GetUser() httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		ctx, span := d.tracer.Start(context.Background(), "get-user")
		defer span.End()

		span.AddEvent("new-get-user-id-request", trace.WithAttributes(attribute.Int("user-count", 87)))
		span.SetAttributes(d.otelUserIDKey.String(params.ByName("id")))

		_, span2 := d.tracer.Start(ctx, "background-gc-job")
		time.Sleep(time.Millisecond * 20)

		d.requestValidator(ctx)
		time.Sleep(time.Millisecond * 30)

		err := d.walletRepository.RegisterToWallet(ctx, params.ByName("id"))
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			_, err = fmt.Fprintf(writer, "err: %+v", err)
			if err != nil {
				log.Println("response print err", err)
			}
			return
		}

		span2.End()

		user, err := d.userRepository.GetUser(ctx, params.ByName("id"))
		if err != nil {
			writer.WriteHeader(http.StatusNotFound)
			_, err = fmt.Fprintf(writer, "err: %+v", err)
			if err != nil {
				log.Println("response print err", err)
			}
			return
		}

		err = d.tellMeYourIPRepository.TellMeYourIP(ctx)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			_, err = fmt.Fprintf(writer, "err: %+v", err)
			if err != nil {
				log.Println("response print err", err)
			}
			return
		}

		err = d.userEventsProducer.Produce(ctx, user.Id)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			_, err = fmt.Fprintf(writer, "err: %+v", err)
			if err != nil {
				log.Println("response print err", err)
			}
			return
		}

		writer.WriteHeader(http.StatusOK)
		_, err = fmt.Fprintf(writer, "default route, time: %s, userName: %s", time.Now(), user.Name)
		if err != nil {
			log.Println("response print err", err)
		}
	}
}

func (d *UserHandler) requestValidator(ctx context.Context) {
	ctx, span := d.tracer.Start(ctx, "request-validator")
	defer span.End()

	time.Sleep(time.Duration(rand.Intn(60)+100) * time.Millisecond)

	span.SetAttributes(d.otelTimeKey.String(time.Now().String()))
	span.AddEvent("successful request validation")
}

type UserRepository interface {
	GetUser(ctx context.Context, id string) (models.User, error)
}

type WalletRepository interface {
	RegisterToWallet(ctx context.Context, id string) error
}

type TellMeYourIPRepository interface {
	TellMeYourIP(ctx context.Context) error
}

type UserEventsProducer interface {
	Produce(ctx context.Context, id string) error
}
