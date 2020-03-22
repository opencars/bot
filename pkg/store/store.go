package store

// Store is responsible for data manipulation.
type Store interface {
	User() UserRepository
	Update() UpdateRepository
}
