package sqlstore

import (
	"context"
	"database/sql"

	"github.com/opencars/bot/pkg/domain/model"
)

type MessageRepository struct {
	s *Store
}

func (r *MessageRepository) Create(ctx context.Context, message *model.Message) error {
	tx, err := r.s.db.BeginTxx(ctx, &sql.TxOptions{Isolation: 0, ReadOnly: false})
	if err != nil {
		return err
	}

	// 1. Create user, if not exist.
	_, err = tx.ExecContext(ctx,
		`INSERT INTO users (
			id, first_name, last_name, username, language_code
		) VALUES (
			$1, $2, $3, $4, $5
		) ON CONFLICT DO NOTHING`,
		message.User.ID,
		message.User.FirstName,
		message.User.LastName,
		message.User.UserName,
		message.User.LanguageCode,
	)

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// 2. Insert message.
	_, err = tx.ExecContext(ctx,
		`INSERT INTO updates (
			id, user_id, text, time
		) VALUES (
			$1, $2, $3, $4
		)`,
		message.ID, message.User.ID, message.Text, message.Time,
	)

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
