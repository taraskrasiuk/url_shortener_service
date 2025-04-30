package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"taraskrasiuk/url_shortener_service/internal/shortener"
	"time"
)

// Middleware for logging request.
// It should display, time, method, url and content type
func ReqInfoMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logTime := start.Format(time.RFC822)
		contentType := r.Header.Get("Content-Type")
		if contentType == "" {
			contentType = r.Header.Get("content-type")
		}
		// server request
		next.ServeHTTP(w, r)

		logMsg := fmt.Sprintf("[LOG]: %s %s Content-Type: %s \t %s [%s]", r.Method, r.URL.Path, contentType, logTime, time.Since(start))
		fmt.Println(logMsg)

		return
	}
}

// POST Handler for creating a shorter link.
// Handler epxects a data in format "multipart/form-data"
// and requirs the filed "link" to exists.
func HandleCreateShortLink(w http.ResponseWriter, r *http.Request) {
	// parse form data
	err := r.ParseMultipartForm(0)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	linkValue := r.FormValue("link")
	if linkValue == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "the link could not being an empty field")
		return
	}

	shortenLink, err := shortener.NewShortLinker(10, "http", "localhost").Create(linkValue)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "could not short a link")
		return
	}
	result := struct {
		ShortenLink string `json:"shortenLink"`
	}{
		shortenLink,
	}
	defer r.Body.Close()

	jsonRes, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "could not make a json")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonRes)
}
