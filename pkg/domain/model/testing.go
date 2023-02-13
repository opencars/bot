package model

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

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

func TestMessage(t *testing.T) *Message {
	t.Helper()

	return &Message{
		ID: rand.Intn(math.MaxInt32),
		User: User{
			ID:        rand.Intn(math.MaxInt32),
			FirstName: "John",
		},
		Text: "/start",
		Time: time.Unix(1584916221, 0).UTC(),
	}
}
