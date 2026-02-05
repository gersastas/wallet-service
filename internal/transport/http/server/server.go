package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	httpServer *http.Server

	mu    sync.Mutex
	stats map[string]int
}

func NewServer(address string) *Server {
	s := &Server{
		stats: make(map[string]int),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/time", s.timeHandler)
	mux.HandleFunc("/stats", s.statsHandler)

	s.httpServer = &http.Server{
		Addr:    address,
		Handler: mux,
	}

	return s
}

func (s *Server) Run() error {
	err := s.httpServer.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) timeHandler(w http.ResponseWriter, r *http.Request) {
	ip := clientIP(r)

	s.mu.Lock()
	s.stats[ip]++
	s.mu.Unlock()

	now := time.Now().Format(time.RFC3339)

	if _, err := w.Write([]byte(now)); err != nil {
		panic(err)
	}
}

func (s *Server) statsHandler(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for ip, count := range s.stats {
		if _, err := fmt.Fprintf(w, "%s: %d\n", ip, count); err != nil {
			panic(err)
		}
	}
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
