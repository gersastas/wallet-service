package server

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func New(address string) *Server {
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

	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now().Format(time.RFC3339)

	if _, err := fmt.Fprint(w, now); err != nil {
		panic(err)
	}
}
