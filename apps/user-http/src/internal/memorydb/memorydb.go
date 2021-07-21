package memorydb

import (
	"context"
	"math/rand"
	"time"

	"user-http/internal/models"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type DB struct {
	tracer        trace.Tracer
	otelUserIDKey attribute.Key
	db            map[string]models.User
}

func New() *DB {
	return &DB{
		tracer:        otel.Tracer("sven.njegac/open-telemetry-k8s"),
		otelUserIDKey: "memory-db/user-id",
		db: map[string]models.User{
			"1x": {
				Id:   "1x",
				Name: "Sven",
				Age:  26,
			},
			"2x": {
				Id:   "2x",
				Name: "Marko",
				Age:  45,
			},
			"3x": {
				Id:   "3x",
				Name: "Ivana",
				Age:  31,
			},
		},
	}
}

func (d *DB) GetUser(ctx context.Context, id string) (models.User, error) {
	ctx, span := d.tracer.Start(ctx, "db-get-user")
	defer span.End()

	span.SetAttributes(d.otelUserIDKey.String(id))

	// fake db delay
	delay := rand.Intn(200)
	delay += 50
	time.Sleep(time.Duration(delay) * time.Millisecond)
	if delay > 170 {
		span.RecordError(errors.New("network record error"))
		span.SetStatus(codes.Error, "")
		span.AddEvent("network failure", trace.WithAttributes(attribute.Int("delay", delay)))
		return models.User{}, errors.New("database network error")
	}

	span.AddEvent("network success", trace.WithAttributes(attribute.Int("delay", delay)))

	user, ok := d.db[id]
	if !ok {
		span.RecordError(
			errors.New("db record error"),
			trace.WithAttributes(
				attribute.String("not-found", "memory-db"),
			),
		)
		span.SetStatus(codes.Error, "")
		span.AddEvent("user not found")
		return models.User{}, errors.New("database not found error")
	}

	d.parsing(ctx)

	span.SetStatus(codes.Ok, "")

	return user, nil
}

func (d *DB) parsing(ctx context.Context) {
	ctx, span := d.tracer.Start(ctx, "parsing")
	defer span.End()

	span.AddEvent("start parsing")

	parseTime := rand.Intn(50) + 50
	time.Sleep(time.Duration(parseTime) * time.Millisecond)

	span.AddEvent("end parsing", trace.WithAttributes(attribute.Int("time spent", parseTime)))
}
