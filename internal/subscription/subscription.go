// Copyright (C) 2019 Ali Shanaakh, github@shanaakh.pro
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package subscription

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type Subscription struct {
	Chat *tgbotapi.Chat
	LastIDs []string
	Quitter chan struct{}
	Running bool
	Params map[string]string
}

func NewSubscription(chat *tgbotapi.Chat, params map[string]string) *Subscription {
	lastIDs := make([]string, 0)
	quitter := make(chan struct{})

	return &Subscription{
		Chat: chat,
		LastIDs: lastIDs,
		Quitter: quitter,
		Running: false,
		Params: params,
	}
}

func (s *Subscription) Stop() {
	// Avoid closing closed channel.
	if !s.Running {
		return
	}

	close(s.Quitter)
	s.Running = false
}

func (s *Subscription) Start(callback func()) {
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
				callback()
			}
		}
	}(s.Quitter)

	// Mark subscription as "Running".
	s.Running = true
}
