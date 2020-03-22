package store

import (
	"github.com/opencars/bot/pkg/model"
)

type UserRepository interface {
	Create(user *model.User) error
	FindByID(id int) (*model.User, error)
}

type UpdateRepository interface {
	Create(update *model.Update) error
	FindByID(id int) (*model.Update, error)
}
