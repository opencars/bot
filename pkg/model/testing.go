package model

import (
	"testing"
	"time"
)

// TestUser returns example of valid User entity.
func TestUser(t *testing.T) *User {
	t.Helper()

	username := "example_user"
	lang := "en"
	return &User{
		ID:           9999999999,
		FirstName:    "Example",
		LastName:     nil,
		UserName:     &username,
		LanguageCode: &lang,
	}
}

// TestUpdate returns example of valid Update entity.
func TestUpdate(t *testing.T) *Update {
	t.Helper()

	return &Update{
		ID:     1,
		UserID: 9999999999,
		Text:   "/start",
		Time:   time.Unix(1584916221, 0).UTC(),
	}
}
