package main

import (
	httpServer "github.com/gersastas/wallet-service/internal/transport/http/server"
	"log"
)

func main() {
	server := httpServer.NewServer(":8080")
	if err := server.Run(); err != nil {
		log.Panic(err)
	}
}
