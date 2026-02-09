package server

import (
	"errors"
	"net"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"
)

type Server struct {
	httpServer *http.Server
	logger     *zap.Logger
	stats      map[string]int
}

func New(address string) *Server {
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

	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}

func (s *Server) timeHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now().Format(time.RFC3339)

	ip := getIP(r)

	s.stats[ip]++
	if _, err := w.Write([]byte(now)); err != nil {
		s.logger.Warn(
			"failed to write response",
			zap.Error(err),
			zap.String("remote_addr", r.RemoteAddr),
			zap.String("path", r.URL.Path),
		)
	}
}

func (s *Server) statsHandler(w http.ResponseWriter, r *http.Request) {
	for ip, count := range s.stats {
		if _, err := w.Write([]byte(ip + "\t" + strconv.Itoa(count) + "\n")); err != nil {
			s.logger.Warn(
				"failed to write response", zap.Error(err))
		}
	}
}

func getIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	return host
}
