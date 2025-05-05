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
		err := st.Drop()
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

func TestFileStoragePersistence(t *testing.T) {
	mockedTestData := map[string]string{
		"295sMlugpK": "https//www.digitalocean.com/community/tutorials/how-to-use-the-flag-package-in-go",
		"5r3-FbHcYD": "https://www.digitalocean.com/community/tutorials/how-to-use-the-flag-package-in-go",
		"pNopWqvEs5": "https://www.digitalocean.com/community/tutorials/how-to-use-the-flag-package-in-go",
		"w1GN_QzF-8": "https://www.digitalocean.com/community/tutorials/how-to-use-the-flag-package-in-go",
	}
	testFileName := "mock_test.db"

	f, _ := os.OpenFile(testFileName, os.O_CREATE|os.O_WRONLY, 0644)
	for k, v := range mockedTestData {
		f.WriteString(k + "=" + v + "\n")
	}
	f.Close()

	st := NewFileStorage(testFileName)
	defer func() {
		err := os.Remove(testFileName)
		if err != nil {
			t.Fatal(err)
		}
	}()

	for k, v := range mockedTestData {
		t.Run("the storage should return a value by a key", func(t *testing.T) {
			foundValue, err := st.Get(k)
			if err != nil {
				t.Fatalf("expected no error but got %v", err)
			}
			if foundValue != v {
				t.Fatalf("expected a returned value for a key %s to be %s ,but got %s", k, v, foundValue)
			}
		})
	}
}
