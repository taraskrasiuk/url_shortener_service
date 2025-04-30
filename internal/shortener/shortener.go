package shortener

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	u "net/url"
	"time"
)

var maxLen = 12

// Validate incoming url string
// 1. Check len of url
// 2. Check the url consintent and can be parsed
// 3. Make a dial request

var dialReq = net.DialTimeout

func validateUrl(url string) error {
	if url == "" {
		return errors.New("url could not being empty")
	}
	u, err := u.ParseRequestURI(url)
	if err != nil {
		return err
	}
	reqDeadline := 500 * time.Millisecond
	conn, err := dialReq("tcp", u.RawPath, reqDeadline)
	// check for a nil, just for testing purposes
	if conn != nil {
		defer conn.Close()
	}
	if err != nil {
		return err
	}
	return nil
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
	l int
}

func NewShortLinker(length int) *ShortLinker {
	return &ShortLinker{length}
}

func (s *ShortLinker) Create(url string) (string, error) {
	err := validateUrl(url)
	if err != nil {
		return "", err
	}

	id, err := generateUUID(s.l)
	if err != nil {
		return "", err
	}
	return id, nil
}
