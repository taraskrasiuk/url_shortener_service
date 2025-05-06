package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"taraskrasiuk/url_shortener_service/internal/shortener"
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

// Handler epxects a data in format "multipart/form-data"
// and requirs the filed "link" to exists.
func (h *UrlShortenerHandler) HandlerCreateShortLink(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprint(w, "incorrect content type")
		return
	}
	// parse form data
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprint(w, err.Error())
		return
	}
	// get a value from a form data, and check it
	linkValue := r.URL.Query().Get("link")
	if linkValue == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
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
