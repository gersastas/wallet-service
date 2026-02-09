package main

import (
	"log"

	httpserver "github.com/gersastas/wallet-service/internal/transport/http/server"
)

func main() {
	server := httpserver.New(":8083")
	if err := server.Run(); err != nil {
		log.Panic(err)
	}
}
