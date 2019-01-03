// Copyright (C) 2019 Ali Shanaakh, github@shanaakh.pro
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package subscription

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"time"
)

type Subscription struct {
	Chat *tgbotapi.Chat
	LastIDs []string
	Quitter chan struct{}
	Running bool
}

func NewSubscription(chat *tgbotapi.Chat) *Subscription {
	lastIDs := make([]string, 0)
	stopper := make(chan struct{})
	return &Subscription{chat, lastIDs, stopper, false}
}


func (s *Subscription) Stop() {
	// Avoid closing closed channel.
	if !s.Running {
		return
	}

	close(s.Quitter)
	s.Running = false
}

func (s *Subscription) Start(b *tgbotapi.BotAPI) {
	// Stop previous goroutine, if it is running.
	if s.Running {
		close(s.Quitter)
		s.Quitter = make(chan struct{})
	}

	go func(quit chan struct{}) {
		for {
			select {
			case <-quit:
				return
			default:
				msg := tgbotapi.NewMessage(s.Chat.ID, "Hello, " + s.Chat.FirstName)

				if _, err := b.Send(msg); err != nil {
					log.Print(err)
				} else {
					log.Printf("Successfully delivered to Chat: %d\n", s.Chat.ID)
				}

				time.Sleep(time.Second * 10)
			}
		}
	}(s.Quitter)

	// Mark subscription as "Running".
	s.Running = true
}
