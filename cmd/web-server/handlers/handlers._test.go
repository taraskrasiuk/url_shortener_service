package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	appconfig "taraskrasiuk/url_shortener_service/cmd/web-server/app_config"
	"taraskrasiuk/url_shortener_service/internal/storage"
	"testing"
)

var cfg = appconfig.NewConfig()

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
		handler := NewUrlShortenerHandler(st, cfg)
		handler.HandlerCreateShortLink(w, req)

		res := w.Result()
		defer res.Body.Close()

		type response struct {
			ShortenLink string `json:"shortenLink"`
		}

		if res.StatusCode != http.StatusBadRequest {
			t.Fatalf("expected status code to be 400 but got %d", res.StatusCode)
		}

		// var resData response
		// b, err := io.ReadAll(res.Body)
		// t.Log("qwe ::" + string(b))
		// if err != nil {
		// 	t.Fatalf("could not read response %v", err)
		// }
		// err = json.Unmarshal(b, &resData)
		// if err != nil {
		// 	t.Fatalf("could not unmarshal response bytes %v", err)
		// }
		// if resData.ShortenLink != fmt.Sprint("qwe") {
		// 	t.Fatalf("expect shorten link to be: %s but got %s", "1", "2")
		// }
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
