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
	"user-http/internal/wallet"

	"github.com/julienschmidt/httprouter"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type DefaultHandler struct {
	tracer                 trace.Tracer
	key1                   attribute.Key
	key2                   attribute.Key
	userRepository         UserRepository
	walletRepo             WalletRepo
	tellMeYourIPRepository TellMeYourIPRepository
}

func NewDefaultHandler() *DefaultHandler {
	return &DefaultHandler{
		tracer:                 otel.Tracer("sven.njegac/basic"),
		key1:                   "sven.njegac/key-1",
		key2:                   "sven.njegac/key-2",
		userRepository:         memorydb.New(),
		walletRepo:             wallet.New(),
		tellMeYourIPRepository: tellip.NewTellIP(),
	}
}

func (d *DefaultHandler) Default() httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		ctx, span := d.tracer.Start(context.Background(), "get-user")
		defer span.End()

		span.AddEvent("Nice operation!", trace.WithAttributes(attribute.Int("bogons", 100)))
		span.SetAttributes(d.key1.String("key-1-yes"))

		_, span2 := d.tracer.Start(ctx, "background-gc-job")
		time.Sleep(time.Millisecond * 20)

		subRes := d.requestValidator(ctx)

		time.Sleep(time.Millisecond * 30)

		err := d.walletRepo.RegisterToWallet(ctx, params.ByName("id"))
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			_, err = fmt.Fprintf(writer, "err: %+v", err)
			if err != nil {
				log.Println("default err", err)
			}
			return
		}

		span2.End()

		user, err := d.userRepository.GetUser(ctx, params.ByName("id"))
		if err != nil {
			writer.WriteHeader(http.StatusNotFound)
			_, err = fmt.Fprintf(writer, "err: %+v", err)
			if err != nil {
				log.Println("default err", err)
			}
			return
		}

		err = d.tellMeYourIPRepository.TellMeYourIP(ctx)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			_, err = fmt.Fprintf(writer, "err: %+v", err)
			if err != nil {
				log.Println("default err", err)
			}
			return
		}

		writer.WriteHeader(http.StatusOK)
		_, err = fmt.Fprintf(writer, "default route, time: %s, subRes: %s, userName: %s", time.Now(), subRes, user.Name)
		if err != nil {
			log.Println("default err", err)
		}
	}
}

func (d *DefaultHandler) requestValidator(ctx context.Context) string {
	var span trace.Span
	ctx, span = d.tracer.Start(ctx, "request-validator")
	defer span.End()

	time.Sleep(time.Duration(rand.Intn(60)+100) * time.Millisecond)

	span.SetAttributes(d.key2.String("key-2-no"))
	span.AddEvent("Sub span event")

	return "sub function result"
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

type UserRepository interface {
	GetUser(ctx context.Context, id string) (models.User, error)
}

type WalletRepo interface {
	RegisterToWallet(ctx context.Context, id string) error
}

type TellMeYourIPRepository interface {
	TellMeYourIP(ctx context.Context) error
}
