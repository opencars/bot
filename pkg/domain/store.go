package domain

// Store is responsible for data manipulation.
type Store interface {
	User() UserRepository
	Update() UpdateRepository
}

type UserRepository interface {
	Create(user *User) error
	FindByID(id int) (*User, error)
}

type UpdateRepository interface {
	Create(update *Update) error
	FindByID(id int) (*Update, error)
}
