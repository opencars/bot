package model

import "time"

type Update struct {
	ID     int       `json:"id" db:"id"`
	UserID int       `json:"user_id" db:"user_id"`
	Text   string    `json:"text" db:"text"`
	Time   time.Time `json:"time" db:"time"`
}
