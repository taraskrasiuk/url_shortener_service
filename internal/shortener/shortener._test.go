package shortener

import (
	"errors"
	"net"
	"net/url"
	httpdialer "taraskrasiuk/url_shortener_service/internal/httpDialer"
	"testing"
)

func TestGenerateUUID(t *testing.T) {
	// Test for failed result
	// change package variable "maxLen" to ten
	maxLen = 10
	// since maxLen is 10, the next generations should fail
	len := 12
	_, err := generateUUID(len)
	if err == nil {
		t.Fatalf("expected to return an error, when len is %d", len)
	}

	len = 0
	_, err = generateUUID(len)
	if err == nil {
		t.Fatalf("expected to return an error, when len is %d", len)
	}

	// Should return uniq ids
	store := make(map[string]struct{})

	for i := 0; i < 10000; i++ {
		id, err := generateUUID(8)
		if err != nil {
			t.Fatalf("expected to generate id, but got an error: %v", err)
		}
		if _, has := store[id]; has {
			t.Fatal("expected no duplicates were generated, but found")
		}
		store[id] = struct{}{}
	}
}

func TestFailedValidateUrlByDialRequest(t *testing.T) {
	// change the function behavior in order to fail
	dialReq = func(_ *url.URL) (net.Conn, error) {
		return nil, errors.New("failed dial request")
	}

	url := "https://wrong-address"

	_, err := parseLink(url)
	if err == nil {
		t.Fatal("expected an error to be occured")
	}
}

func TestFailedValidateUrlInvalidURI(t *testing.T) {
	invalidUri := "invalid_uri"
	_, err := parseLink(invalidUri)
	if err == nil {
		t.Fatal("expected an error to be occured")
	}
}

func TestCreateShorterLink(t *testing.T) {
	// Revert mocked function
	dialReq = httpdialer.HttpDialProxy

	inputLink := "https://www.digitalocean.com/community/tutorials/how-to-use-the-flag-package-in-go"

	s := NewShortLinker(10, "http", "localhost")

	shorterLink, err := s.Create(inputLink)
	if err != nil {
		t.Errorf("expected to get a shorter link, but got an error: %v", err)
	}
	shortUrl, err := url.Parse(shorterLink)
	if err != nil {
		t.Errorf("expected to parse a shorter link, but got an error: %v", err)
	}
	if shortUrl.Scheme != "http" || shortUrl.Host != "localhost" {
		t.Error("shorter link does not have a provided scheme and host from env")
	}
}
