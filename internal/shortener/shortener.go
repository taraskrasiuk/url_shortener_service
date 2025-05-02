package shortener

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	httpdialer "taraskrasiuk/url_shortener_service/internal/httpDialer"
)

var maxLen = 12

// Validate incoming url string
// 1. Check len of url
// 2. Check the url consintent and can be parsed
// 3. Make a dial request
var dialReq = httpdialer.HttpDialProxy

func parseLink(link string) (*url.URL, error) {
	if link == "" {
		return nil, errors.New("url could not being empty")
	}
	u, err := url.Parse(link)
	if err != nil {
		return nil, err
	}

	conn, err := dialReq(u)
	// check for a nil, just for testing purposes
	if conn != nil {
		defer conn.Close()
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

func generateUUID(length int) (string, error) {
	if length <= 0 || length > maxLen {
		return "", errors.New(fmt.Sprintf("generateUUID: expected len to be > 0 and < %d", maxLen))
	}
	// create a slice of bytes
	buf := make([]byte, length)

	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}
	// encode bytes to string
	str := base64.RawURLEncoding.EncodeToString(buf)
	// ensure in case of long string, it will be cutted
	if length < len(str) {
		str = str[:length]
	}

	return str, nil
}

type ShortLinker struct {
	len int
}

func NewShortLinker(length int) *ShortLinker {
	return &ShortLinker{length}
}

func (s *ShortLinker) Create(link string) (string, error) {
	_, err := parseLink(link)
	if err != nil {
		return "", err
	}

	id, err := generateUUID(s.len)
	if err != nil {
		return "", err
	}
	return id, nil
}
