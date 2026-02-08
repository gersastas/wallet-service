package main

import (
	httpserver "github.com/gersastas/wallet-service/internal/transport/http/server"
	"log"
)

func main() {
	server := httpserver.New(":8081")
	if err := server.Run(); err != nil {
		log.Panic(err)
	}
}
