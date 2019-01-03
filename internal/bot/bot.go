// Copyright (C) 2019 Ali Shanaakh, github@shanaakh.pro
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package bot

import (
	"fmt"
	"log"
	"strings"
	"time"
)

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/shal/robot/pkg/util"
)

var (
	listeners map[int64]chan struct{}
)

func newBot() *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(util.MustGetEnv("BOT_TOKEN"))

	if err != nil {
		panic(err)
	}

	bot.Debug = true
	log.Printf("Bot authorized %s", bot.Self.UserName)

	return bot
}

func processUpdate(u tgbotapi.Update, b *tgbotapi.BotAPI) {
	if u.Message != nil {
		processMsgUpdate(u.Message, b)
	}
}

func processMsgUpdate(m *tgbotapi.Message, b *tgbotapi.BotAPI) {
	log.Printf("Recieved message from [%d] %s\n", m.Chat.ID, m.Chat.UserName)
	log.Println(m.Text)

	if isFollowCmd(m.Text) {
		processFollowMsg(m, b)
	} else if isStopCmd(m.Text) {
		close(listeners[m.Chat.ID])
	}
}

func isFollowCmd(lexeme string) bool {
	return strings.HasPrefix(lexeme, "/follow")
}

func isStopCmd(lexeme string) bool {
	return strings.HasPrefix(lexeme, "/stop")
}

func processFollowMsg(m *tgbotapi.Message, b *tgbotapi.BotAPI) {
	//lexemes := strings.Split(m.Text, " ")

	//if len(lexemes) < 2 {
	//	return
	//} else if !strings.HasPrefix(lexemes[1], "https://auto.ria.com/search") {
	//	return
	//}

	//autoRia := autoria_api.NewAPI(util.MustGetEnv("RIA_API_KEY"))
	//params, err := autoria_api.ParseCarSearchParams(lexemes[1])

	//if err != nil {
	//	return
	//}
	//
	//search := autoRia.GetSearchCars(params...)
	//
	//msg := tgbotapi.NewMessage(m.Chat.ID, "Hi")
	//msg.Text = strings.Join(search.Result.SearchResult.CarsIDs, " ")
	//
	//if _, err = b.Send(msg); err != nil {
	//	fmt.Print(err)
	//}

	registerListener(m.Chat.ID)
}

func registerListener(ID int64) {
	listeners[ID] = make(chan struct{})

	go func(q chan struct{}) {
		for {
			select {
			case <-q:
				return
			default:
				fmt.Print("Hi --------> ", ID)
				time.Sleep(time.Second * 3)
			}
		}
	}(listeners[ID])
}

func Run() {
	listeners = make(map[int64]chan struct{})
	bot := newBot()

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates, err := bot.GetUpdatesChan(updateConfig)

	if err != nil {
		log.Println(err.Error())
	}

	for update := range updates {
		processUpdate(update, bot)
	}
}
