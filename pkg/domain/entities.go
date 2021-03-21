package domain

import "time"

type User struct {
	ID           int
	FirstName    string
	LastName     *string
	UserName     *string
	LanguageCode *string
}

type Message struct {
	ID   int
	User User
	Text string
	Time time.Time
}
