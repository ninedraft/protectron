package registry

import (
	"github.com/asdine/storm"
)

type Registry struct {
	db *storm.DB
}

func New(dbPath string) (Registry, error) {
	var stDB, err = storm.Open(dbPath)
	if err != nil {
		return Registry{}, err
	}
	return Registry{
		db: stDB,
	}, nil
}
