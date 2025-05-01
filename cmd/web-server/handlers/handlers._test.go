package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"taraskrasiuk/url_shortener_service/internal/storage"
	"testing"
)

func TestSuccessHandlerCreateShortLink(t *testing.T) {
	st := storage.NewFileStorage("test.db")
	defer st.Drop()

	buf := &bytes.Buffer{}

	writer := multipart.NewWriter(buf)
	writer.WriteField("link", "https://www.digitalocean.com/community/tutorials/how-to-use-the-flag-package-in-go")
	// close the writer
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/shorten", buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	// run handler
	handler := NewUrlShortenerHandler(st)
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
	if err := json.Unmarshal(data, respContent); err != nil {
		t.Fatalf("could not unmarshal the response data %v", err)
	}
}

func TestFailHandlerCreateShortLink(t *testing.T) {
	tests := []struct {
		body        [2]string
		contetnType string
	}{
		{
			body:        [2]string{"incorrect_field_name", "https://google.com"}, // incorrect field name
			contetnType: "multipart/form-data",
		},
		{

			body:        [2]string{"link", ""}, // missed link
			contetnType: "multipart/form-data",
		},
		{

			body:        [2]string{"link", "https://google.com"},
			contetnType: "application/json", // incorrect content type
		},
	}

	st := storage.NewFileStorage("test.db")
	defer st.Drop()

	for _, test := range tests {
		buf := &bytes.Buffer{}

		writer := multipart.NewWriter(buf)
		writer.WriteField(test.body[0], test.body[1])
		// close the writer
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/shorten", buf)
		req.Header.Set("Content-Type", test.contetnType)

		w := httptest.NewRecorder()
		// run handler
		handler := NewUrlShortenerHandler(st)
		handler.HandlerCreateShortLink(w, req)

		res := w.Result()
		defer res.Body.Close()

		type response struct {
			ShortenLink string `json:"shortenLink"`
		}

		if res.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected status code to be 400 but got %d", res.StatusCode)
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

	handler := NewUrlShortenerHandler(st)
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
