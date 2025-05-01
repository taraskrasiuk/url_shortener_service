package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"taraskrasiuk/url_shortener_service/cmd/web-server/handlers"
	"taraskrasiuk/url_shortener_service/internal/storage"
)

func main() {
	host := "localhost"
	port := "8080"

	mux := http.NewServeMux()

	var appStorage storage.Storage
	appStorage = storage.NewFileStorage("url_storage.db")
	urlShortenerHandler := handlers.NewUrlShortenerHandler(appStorage)

	mux.HandleFunc("POST /shorten", urlShortenerHandler.HandlerCreateShortLink)
	mux.HandleFunc("GET /{shortenID}", urlShortenerHandler.HandleShortLink)

	// register middlewares
	handler := handlers.ReqInfoMiddleware(mux)

	addr := fmt.Sprintf("%s:%s", host, port)

	fmt.Println("Server is running: " + addr)
	if err := http.ListenAndServe(addr, handler); err != nil && errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
