package sqlstore

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/opencars/bot/pkg/config"
	"github.com/opencars/bot/pkg/domain"
)

type Store struct {
	db *sqlx.DB

	messageRepository *MessageRepository
}

func (s *Store) Message() domain.MessageRepository {
	if s.messageRepository == nil {
		s.messageRepository = &MessageRepository{
			s: s,
		}
	}

	return s.messageRepository
}

func New(conf *config.Database) (*Store, error) {
	dataSourceName := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s password=%s",
		conf.Host, conf.Port, conf.User, conf.Database, conf.SSLMode, conf.Password,
	)

	db, err := sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	return &Store{
		db: db,
	}, nil
}
