package storage

import (
	"os"
	"testing"
)

func TestWriteFileStorage(t *testing.T) {
	testFileName := "mock_test.db"
	st := NewFileStorage(testFileName)
	defer func() {
		err := os.Remove(testFileName)
		if err != nil {
			t.Fatal(err)
		}
	}()
	defer st.f.Close()

	pairs := map[string]string{
		"https://google.com": "http://local?123",
		"https://yahoo.com":  "http://local?qweASD",
	}

	for k, v := range pairs {
		err := st.Write(k, v)
		if err != nil {
			t.Fatalf("not expected on error: %v", err)
		}
	}

	for k, v := range pairs {
		val, err := st.Get(k)
		if err != nil {
			t.Fatalf("not expected on error: %v", err)
		}
		if val != v {
			t.Fatalf("values do no mathc, got %s but want %s", val, v)
		}
	}
}
