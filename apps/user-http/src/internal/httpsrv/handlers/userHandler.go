package handlers

import (
	"context"
	"fmt"
	"log"
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
	userRepository         UserRepository
	walletRepository       WalletRepository
	tellMeYourIPRepository TellMeYourIPRepository
	userEventsProducer     UserEventsProducer
}

func NewDefaultHandler() *UserHandler {
	return &UserHandler{
		tracer:                 otel.Tracer("sven.njegac/open-telemetry-k8s"),
		userRepository:         memorydb.New(),
		walletRepository:       wallet.New(),
		tellMeYourIPRepository: tellip.NewTellIP(),
		userEventsProducer:     userevents.NewProducer(),
	}
}

func (d *UserHandler) GetUser() httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		ctx, span := otel.Tracer("sven.njegac/open-telemetry-k8s").Start(context.Background(), "get-user")
		defer span.End()

		// Adds new log (event).
		span.AddEvent(
			"new-get-user-id-request",
			trace.WithAttributes(
				attribute.Int("user-count", 87),
				attribute.String("user-id", params.ByName("id")),
			),
			trace.WithTimestamp(time.Now()),
			trace.WithStackTrace(true),
		)
		// Set tag on the span.
		span.SetAttributes(attribute.String("user-id", params.ByName("id")))

		// Fake delay.
		time.Sleep(time.Millisecond * 30)

		// Register user to wallet. (HTTP call)
		err := d.walletRepository.RegisterToWallet(ctx, params.ByName("id"))
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			_, err = fmt.Fprintf(writer, "err: %+v", err)
			if err != nil {
				log.Println("wallet err", err)
			}
			return
		}

		// Get user from user collection. (in memory DB)
		user, err := d.userRepository.GetUser(ctx, params.ByName("id"))
		if err != nil {
			writer.WriteHeader(http.StatusNotFound)
			_, err = fmt.Fprintf(writer, "err: %+v", err)
			if err != nil {
				log.Println("user err", err)
			}
			return
		}

		// Get IP of
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
