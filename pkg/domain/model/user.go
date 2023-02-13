package model

type User struct {
	ID           int
	FirstName    string
	LastName     *string
	UserName     *string
	LanguageCode *string
}
