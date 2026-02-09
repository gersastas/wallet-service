package server

import (
	"errors"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Server struct {
	httpServer *http.Server
	logger     *zap.Logger
}

func New(address string) *Server {
	s := &Server{}

	mux := http.NewServeMux()
	mux.HandleFunc("/time", s.timeHandler)

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

func (s *Server) timeHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now().Format(time.RFC3339)

	if _, err := w.Write([]byte(now)); err != nil {
		s.logger.Warn(
			"failed to write response",
			zap.Error(err),
			zap.String("remote_addr", r.RemoteAddr),
			zap.String("path", r.URL.Path),
		)
	}
}
