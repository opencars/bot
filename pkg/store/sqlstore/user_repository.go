package sqlstore

import (
	"github.com/opencars/bot/pkg/model"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Create(user *model.User) error {
	_, err := r.store.db.NamedExec(
		`INSERT INTO users (
			id, first_name, last_name, username, language_code
		) VALUES (
			:id, :first_name, :last_name, :username, :language_code
		) ON CONFLICT DO NOTHING`,
		user,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) FindByID(id int) (*model.User, error) {
	var user model.User

	err := r.store.db.Get(&user,
		`SELECT id, first_name, last_name, username, language_code 
		FROM users WHERE id = $1`,
		id,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
