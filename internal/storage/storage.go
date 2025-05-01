package storage

type Storage interface {
	Write(k, v string) error
	Get(k string) (string, error)
	Drop() error
}
