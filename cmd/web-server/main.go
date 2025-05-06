package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	appconfig "taraskrasiuk/url_shortener_service/cmd/web-server/envConfig"
	"taraskrasiuk/url_shortener_service/cmd/web-server/handlers"
	"taraskrasiuk/url_shortener_service/cmd/web-server/middlewares"
	"taraskrasiuk/url_shortener_service/internal/storage"
)

func main() {
	mux := http.NewServeMux()
	var cfg = appconfig.NewEnvConfig()
	var appStorage = storage.NewFileStorage(os.Getenv("STORAGE_FILE_PATH"))

	urlShortenerHandler := handlers.NewUrlShortenerHandler(appStorage, cfg)

	mux.HandleFunc("POST /shorten", urlShortenerHandler.HandlerCreateShortLink)
	mux.HandleFunc("GET /{shortenID}", urlShortenerHandler.HandleShortLink)

	// register middlewares
	handler := middlewares.NewLoggerMiddleware(mux, os.Stdout)
	withRPSLimiter := middlewares.NewRateLimiterMiddleware(handler, 100)

	addr := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))

	fmt.Println("Server is running: " + addr)
	if err := http.ListenAndServe(addr, withRPSLimiter); err != nil && errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
