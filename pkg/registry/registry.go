package registry

import (
	"github.com/ninedraft/protectron/pkg/user"

	"github.com/asdine/storm"
)

type Registry struct {
	db *storm.DB
}

func New(dbPath string) (Registry, error) {
	var stDB, err = storm.Open(dbPath, defaultStormOptions())
	if err != nil {
		return Registry{}, err
	}
	return Registry{
		db: stDB,
	}, nil
}

func (registry Registry) AddUser(u user.User) error {
	var tx, err = registry.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := registry.db.Save(&u); err != nil || err != storm.ErrAlreadyExists {
		return err
	}
	return tx.Commit()
}

func (registry Registry) IsBanned(telegramID int64) (bool, error) {
	var tx, err = registry.db.Begin(false)
	if err != nil {
		return false, err
	}
	var u user.User
	if err := tx.One("TelegramID", telegramID, &u); err != nil {
		return false, err
	}
	return u.IsBanned, tx.Commit()
}

func (registry Registry) BanUser(telegramID int64) error {
	var tx, err = registry.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := tx.UpdateField(&user.User{
		TelegramID: telegramID,
	}, "IsBanned", true); err != nil {
		return err
	}
	return tx.Commit()
}
