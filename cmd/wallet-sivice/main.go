package main

import (
	"log"

	"wallet-service/internal/transport/http/server"
)

func main() {
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
