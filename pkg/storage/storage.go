package storage

type Storage interface {
	Write([]byte) error
}
