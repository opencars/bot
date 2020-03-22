package sqlstore

import (
	"github.com/opencars/bot/pkg/model"
)

type UpdateRepository struct {
	store *Store
}

func (r *UpdateRepository) Create(update *model.Update) error {
	_, err := r.store.db.NamedExec(
		`INSERT INTO updates (
			id, user_id, text, time
		) VALUES (
			:id, :user_id, :text, :time
		)`,
		update,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *UpdateRepository) FindByID(id int) (*model.Update, error) {
	var update model.Update

	err := r.store.db.Get(&update,
		`SELECT id, user_id, text, time
		FROM updates WHERE id = $1`,
		id,
	)

	if err != nil {
		return nil, err
	}

	update.Time = update.Time.UTC()
	return &update, nil
}
