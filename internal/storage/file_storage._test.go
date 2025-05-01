package storage

import (
	"os"
	"sync"
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

func TestFileStorageConcurrency(t *testing.T) {
	testFileName := "mock_test.db"
	st := NewFileStorage(testFileName)
	defer func() {
		err := os.Remove(testFileName)
		if err != nil {
			t.Fatal(err)
		}
	}()

	tests := map[string]string{
		"url-1": "shorten-url-1",
		"url-2": "shorten-url-2",
	}
	var wg sync.WaitGroup
	wg.Add(len(tests))

	for k, v := range tests {
		go func() {
			defer wg.Done()
			st.Write(k, v)
		}()
	}

	wg.Wait()
	for k, v := range tests {
		if storedVal, err := st.Get(k); err != nil {
			t.Fatalf("expected no error but got %v", err)
		} else {
			if storedVal != v {
				t.Fatal("expected a value to be written")
			}
		}
	}
}
