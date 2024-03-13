package store

type store struct {
}

type Store interface {
}

func New() Store {
	return &store{}
}
