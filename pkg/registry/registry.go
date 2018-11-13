package registry

import (
	"github.com/ninedraft/protectron/pkg/user"

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

func (registry Registry) IsBanned(telegramID int64) (bool, error) {
	var u user.User
	if err := registry.db.One("TelegramID", telegramID, &u); err != nil {
		return false, err
	}
	return u.IsBanned, nil
}
