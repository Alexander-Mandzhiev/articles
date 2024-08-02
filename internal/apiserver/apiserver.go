package apiserver

import (
	"context"
	"net/http"
	"time"
)

type APIServer struct {
	httpserver *http.Server
}

func (s *APIServer) Start(port string, handler http.Handler) error {

	s.httpserver = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
	}
	return s.httpserver.ListenAndServe()
}

func (s *APIServer) Shutdown(ctx context.Context) error {
	return s.httpserver.Shutdown(ctx)
}
