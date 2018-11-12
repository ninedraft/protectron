package registry

import (
	"github.com/asdine/storm"
)

type Registry struct {
	db *storm.DB
}

func New() (Registry, error) {
	return Registry{}, nil
}
