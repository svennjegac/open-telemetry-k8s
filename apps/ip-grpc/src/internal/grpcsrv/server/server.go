package server

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"ip-grpc/internal/ip"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
)

const (
	port = 8113
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Start(ctx context.Context) {
	// create grpc server and register ip service handler
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     time.Duration(3600000) * time.Millisecond,
			MaxConnectionAge:      time.Duration(12000) * time.Millisecond,
			MaxConnectionAgeGrace: time.Duration(5000) * time.Millisecond,
			Time:                  time.Duration(3600000) * time.Millisecond,
			Timeout:               time.Duration(20000) * time.Millisecond,
		}),
	)

	ipSvc := &ipService{
		ip:     "99.66.33.11",
		tracer: otel.Tracer("sven.njegac/basic-2"),
	}
	ip.RegisterIPServiceServer(grpcServer, ipSvc)

	// initialize listener for incoming tcp connections
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("error listening for tcp connections; port=%d, err=%s\n", port, err)
	}
	defer lis.Close()

	// start listening for grpc requests
	log.Printf("ip service started; port=%d\n", port)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("error serving grpc requests; err=%s\n", err)
	}
}

type ipService struct {
	ip     string
	tracer trace.Tracer
}

func (i *ipService) TellMeYourIP(ctx context.Context, req *ip.TellMeYourIPRequest) (*ip.TellMeYourIPResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "nil request received")
	}

	ctx, span := i.tracer.Start(ctx, "tell-me-your-ip")
	defer span.End()

	log.Println("handling request;", req.ClientIp)
	time.Sleep(time.Duration(rand.Intn(50)+30) * time.Millisecond)

	return &ip.TellMeYourIPResponse{
		ServerIp: i.ip,
	}, nil
}
