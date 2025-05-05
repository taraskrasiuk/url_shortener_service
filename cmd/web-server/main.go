package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	appconfig "taraskrasiuk/url_shortener_service/cmd/web-server/app_config"
	"taraskrasiuk/url_shortener_service/cmd/web-server/handlers"
	"taraskrasiuk/url_shortener_service/internal/storage"

	"github.com/joho/godotenv"
)

type Config struct {
	scheme string
	host   string
}

func NewConfig() *Config {
	godotenv.Load()

	scheme := os.Getenv("scheme")
	if scheme == "" {
		scheme = "http"
	}
	host := os.Getenv("host")
	if host == "" {
		host = "localhost"
	}
	return &Config{scheme, host}
}

func main() {
	host := "localhost"
	port := "8080"

	mux := http.NewServeMux()

	var appStorage = storage.NewFileStorage("url_storage.db")
	var cfg = appconfig.NewConfig()
	urlShortenerHandler := handlers.NewUrlShortenerHandler(appStorage, cfg)

	mux.HandleFunc("POST /shorten", urlShortenerHandler.HandlerCreateShortLink)
	mux.HandleFunc("GET /{shortenID}", urlShortenerHandler.HandleShortLink)

	// register middlewares
	handler := handlers.ReqInfoMiddleware(mux)
	// handler = middlewares.ReqRateLimit(mux)
	addr := fmt.Sprintf("%s:%s", host, port)

	fmt.Println("Server is running: " + addr)
	if err := http.ListenAndServe(addr, handler); err != nil && errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
