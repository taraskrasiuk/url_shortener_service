package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"taraskrasiuk/url_shortener_service/internal/shortener"
	"time"
)

type dbStorage interface {
	Write(k, v string) error
	Get(k string) (string, error)
	Drop() error
}

type config interface {
	GetScheme() string
	GetHost() string
}

type UrlShortenerHandler struct {
	storage dbStorage
	config  config
}

func NewUrlShortenerHandler(s dbStorage, c config) *UrlShortenerHandler {
	return &UrlShortenerHandler{s, c}
}

// Middleware for logging request.
// It should display, time, method, url and content type
//
// TODO: move to another file
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
func (h *UrlShortenerHandler) HandlerCreateShortLink(w http.ResponseWriter, r *http.Request) {
	// parse form data
	err := r.ParseMultipartForm(0)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}
	// get a value from a form data, and check it
	linkValue := r.FormValue("link")
	if linkValue == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "the link could not being an empty field")
		return
	}

	// create a short version of the link
	shortenID, err := shortener.NewShortLinker(10).Create(linkValue)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "could not short a link")
		return
	}
	// save the short and original version to storage
	err = h.storage.Write(shortenID, linkValue)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "could not save to storage")
		return
	}

	result := struct {
		ShortenLink string `json:"shortenLink"`
	}{
		ShortenLink: fmt.Sprintf("%s://%s/%s", h.config.GetScheme(), h.config.GetHost(), shortenID),
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

// Handler which going to count and gather the information
// and then redirect user to original url.
func (h *UrlShortenerHandler) HandleShortLink(w http.ResponseWriter, r *http.Request) {
	shortId := r.PathValue("shortenID")
	defer r.Body.Close()
	// validate a shortenID
	if strings.TrimSpace(shortId) == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "short id is missed")
		return
	}
	origLink, err := h.storage.Get(shortId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if origLink == "" {
			fmt.Fprint(w, "the short link does not exist")
		} else {
			fmt.Fprint(w, "could not find a link")
		}
		return
	}

	// redirect user to original link
	http.Redirect(w, r, origLink, http.StatusPermanentRedirect)
}
