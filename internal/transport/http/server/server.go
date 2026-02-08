package server

import (
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(address string) *Server {
	s := &Server{}

	mux := http.NewServeMux()
	mux.HandleFunc("/time", timeHandler)

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

func timeHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now().Format(time.RFC3339)
	fmt.Fprint(w, now)
}
