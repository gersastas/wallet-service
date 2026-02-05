package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpserver "github.com/gersastas/wallet-service/internal/transport/http/server"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/time", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(time.Now().Format(time.RFC3339)))
	})

	server := httpserver.NewServer(":8080", mux)

	go func() {
		if err := server.Run(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	}
}
