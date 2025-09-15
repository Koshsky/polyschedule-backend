package httpserver

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	addr   string
	server *http.Server
}

func New(addr string) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	return &Server{addr: addr, server: &http.Server{Addr: addr, Handler: mux}}
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}

// Server returns the internal *http.Server (for external HTTP stack integration)
func (s *Server) Server() *http.Server {
	return s.server
}
