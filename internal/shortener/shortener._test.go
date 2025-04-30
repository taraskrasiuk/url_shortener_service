package shortener

import (
	"errors"
	"net"
	"testing"
	"time"
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
	dialReq = func(_, _ string, _ time.Duration) (net.Conn, error) {
		return nil, errors.New("failed dial request")
	}

	url := "https://wrong-address"

	err := validateUrl(url)
	if err == nil {
		t.Fatal("expected an error to be occured")
	}
}
