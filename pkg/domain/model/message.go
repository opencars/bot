package model

import (
	"time"
)

type Message struct {
	ID   int
	User User
	Text string
	Time time.Time
}
