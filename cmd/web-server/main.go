package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"taraskrasiuk/url_shortener_service/cmd/web-server/handlers"
)

func main() {
	host := "localhost"
	port := "8080"

	mux := http.NewServeMux()

	mux.HandleFunc("POST /shorten", handlers.HandlerCreateShortLink)

	// register middlewares
	handler := handlers.ReqInfoMiddleware(mux)

	addr := fmt.Sprintf("%s:%s", host, port)

	fmt.Println("Server is running: " + addr)
	if err := http.ListenAndServe(addr, handler); err != nil && errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
