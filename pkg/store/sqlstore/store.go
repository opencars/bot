package sqlstore

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/opencars/bot/pkg/config"
	"github.com/opencars/bot/pkg/store"
)

type Store struct {
	db *sqlx.DB

	userRepository   *UserRepository
	updateRepository *UpdateRepository
}

func (s *Store) User() store.UserRepository {
	if s.userRepository == nil {
		s.userRepository = &UserRepository{
			store: s,
		}
	}

	return s.userRepository
}

func (s *Store) Update() store.UpdateRepository {
	if s.updateRepository == nil {
		s.updateRepository = &UpdateRepository{
			store: s,
		}
	}

	return s.updateRepository
}

func New(conf *config.Store) (*Store, error) {
	dataSourceName := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s",
		conf.Host, conf.Port, conf.User, conf.Database, conf.Password,
	)

	db, err := sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	return &Store{
		db: db,
	}, nil
}
