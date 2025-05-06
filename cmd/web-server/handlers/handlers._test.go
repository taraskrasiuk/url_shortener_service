package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	appconfig "taraskrasiuk/url_shortener_service/cmd/web-server/envConfig"
	"taraskrasiuk/url_shortener_service/internal/storage"
	"testing"
)

var cfg = appconfig.NewEnvConfig()

func TestSuccessHandlerCreateShortLink(t *testing.T) {
	st := storage.NewFileStorage("test.db")
	defer st.Drop()

	testLink := "https://www.digitalocean.com/community/tutorials/how-to-use-the-flag-package-in-go"
	url := fmt.Sprintf("/shorten?link=%s", testLink)
	req := httptest.NewRequest(http.MethodPost, url, nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	// run handler
	handler := NewUrlShortenerHandler(st, cfg)
	handler.HandlerCreateShortLink(w, req)

	res := w.Result()
	defer res.Body.Close()

	type response struct {
		ShortenLink string `json:"shortenLink"`
	}
	respContent := &response{}
	data, err := io.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		t.Fatal("the response code should be 200")
	}
	if err != nil {
		t.Fatal("could not read bytes from the response body")
	}
	// The response should containe a json
	if err := json.Unmarshal(data, respContent); err != nil {
		t.Fatalf("could not unmarshal the response data %v", err)
	} else {
		if respContent.ShortenLink == "" {
			t.Fatalf("expected a short id to be not empty")
		}
	}
}

func TestFailHandlerCreateShortLink(t *testing.T) {
	tests := []struct {
		link        string
		contetnType string
	}{
		{

			link:        "", // missed link
			contetnType: "application/x-www-form-urlencoded",
		},
		{

			link:        "https://google.com",
			contetnType: "application/json", // incorrect content type
		},
	}

	st := storage.NewFileStorage("test.db")
	defer st.Drop()

	for _, test := range tests {
		buf := &bytes.Buffer{}

		req := httptest.NewRequest(http.MethodPost, "/shorten"+"?link="+test.link, buf)
		req.Header.Set("Content-Type", test.contetnType)

		w := httptest.NewRecorder()
		// run handler
		handler := NewUrlShortenerHandler(st, cfg)
		handler.HandlerCreateShortLink(w, req)

		res := w.Result()
		defer res.Body.Close()

		type response struct {
			ShortenLink string `json:"shortenLink"`
		}

		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Fatalf("expected status code to be 422 but got %d", res.StatusCode)
		}
	}
}

func TestHandleShortLink(t *testing.T) {
	st := storage.NewFileStorage("test.db")
	defer st.Drop()

	shortLink := "qwe123poi"
	origLink := "http://google.com?longitem=1&and=2"

	err := st.Write(shortLink, origLink)
	if err != nil {
		t.Fatal(err)
		t.FailNow()
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.SetPathValue("shortenID", shortLink)
	w := httptest.NewRecorder()

	handler := NewUrlShortenerHandler(st, cfg)
	handler.HandleShortLink(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusPermanentRedirect {
		t.Fatalf("expect the status code to be %d but got %d", http.StatusPermanentRedirect, res.StatusCode)
	}
	targetUrl := res.Header.Get("Location")
	if targetUrl != origLink {
		t.Fatalf("expect res to contain redirection link to be %s but got %s", origLink, targetUrl)
	}
}
