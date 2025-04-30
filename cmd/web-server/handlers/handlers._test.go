package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSuccessHandleCreateShortLink(t *testing.T) {
	buf := &bytes.Buffer{}

	writer := multipart.NewWriter(buf)
	writer.WriteField("link", "https://www.digitalocean.com/community/tutorials/how-to-use-the-flag-package-in-go")
	// close the writer
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/shorten", buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	// run handler
	HandleCreateShortLink(w, req)

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
