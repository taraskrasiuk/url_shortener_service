package storage

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

type FileStorage struct {
	f      *os.File
	values map[string]string
	mu     sync.Mutex
}

// Create file storage function which returns an instance of FileStorage.
// If file for a storage does not exists, it will create it.
func NewFileStorage(filePath string) *FileStorage {
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fileStorage := &FileStorage{
		f:      f,
		values: make(map[string]string),
	}
	err = fileStorage.readFileStorage()
	if err != nil {
		log.Fatal(err)
	}
	return fileStorage
}

// Read the file and fill the values map in instance.
func (s *FileStorage) readFileStorage() error {
	for {
		data, _, err := bufio.NewReader(s.f).ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		var (
			row = strings.Split(string(data), "=")
			k   = row[0]
			v   = row[1]
		)
		s.values[k] = v
	}
	return nil
}

func (s *FileStorage) Write(key, val string) error {
	if strings.TrimSpace(key) == "" || strings.TrimSpace(val) == "" {
		return errors.New("fileStorage: missed key or val")
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if _, has := s.values[key]; has {
		return errors.New("fileStorage: key duplicated " + key)
	}

	_, err := s.f.Write([]byte(fmt.Sprintf("%s=%s\n", key, val)))
	if err != nil {
		return err
	}

	s.values[key] = val

	return nil
}

func (s *FileStorage) Get(key string) (string, error) {
	if strings.TrimSpace(key) == "" {
		return "", errors.New("fileStorage: missed key")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	return s.values[key], nil
}

func (s *FileStorage) Drop() error {
	defer os.Remove(s.f.Name())
	defer s.f.Close()
	return nil
}
