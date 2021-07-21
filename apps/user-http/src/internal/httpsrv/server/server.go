package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

type Server struct {
	httpServer  *http.Server
	userHandler UserHandler
}

func New(
	userHandler UserHandler,
) *Server {
	s := &Server{
		httpServer: &http.Server{
			Addr:              ":8111",
			ReadTimeout:       time.Second * 3,
			ReadHeaderTimeout: time.Second * 3,
			WriteTimeout:      time.Second * 10,
			IdleTimeout:       time.Second * 60,
		},
		userHandler: userHandler,
	}
	s.setRoutes()
	return s
}

func (s *Server) Start(ctx context.Context) {
	serverErrs := make(chan error, 1)
	defer close(serverErrs)

	go func() {
		log.Println("starting http server")
		err := s.httpServer.ListenAndServe()
		serverErrs <- err
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		err := s.httpServer.Shutdown(shutdownCtx)
		if err != nil {
			log.Println("server shutdown error", err)
		} else {
			log.Println("server stopped gracefully")
		}

	case err := <-serverErrs:
		if err != http.ErrServerClosed {
			log.Println("server unexpected error", err)

			shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			defer cancel()

			err = s.httpServer.Shutdown(shutdownCtx)
			if err != nil {
				log.Println("server shutdown error", err)
			} else {
				log.Println("server stopped gracefully")
			}
		} else {
			log.Println("server listen and serve finished with closed error", err)
		}
	}
}

type UserHandler interface {
	GetUser() httprouter.Handle
}
