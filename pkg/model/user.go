package model

// User represents bot user entity.
type User struct {
	ID           int     `json:"id" db:"id"`
	FirstName    string  `json:"first_name" db:"first_name"`
	LastName     *string `json:"last_name" db:"last_name"`
	UserName     *string `json:"username" db:"username"`
	LanguageCode *string `json:"language_code" db:"language_code"`
}
