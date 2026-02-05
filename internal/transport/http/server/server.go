package server

import (
	"fmt"
	"net/http"
	"time"
)

func Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/time", timeHandler)

	return http.ListenAndServe(":8081", mux)
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now().Format(time.RFC3339)
	fmt.Fprint(w, now)
}
