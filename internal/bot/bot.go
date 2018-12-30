// Copyright (C) 2019 Ali Shanaakh, github@shanaakh.pro
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package bot

import (
	"github.com/shal/robot/pkg/autoria-api"
	"log"
	"strings"
)

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/shal/robot/pkg/util"
)

func Run() {
	token := util.MustGetEnv("BOT_TOKEN")
	apiKey := util.MustGetEnv("RIA_API_KEY")

	bot, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates, err := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		msg := update.Message

		if msg == nil {
			continue
		}

		log.Printf("Recieved message from [%d] %s\n", msg.Chat.ID, msg.Chat.UserName)
		log.Println(msg.Text)

		if strings.HasPrefix(msg.Text, "/follow") {
			lexemes := strings.Split(msg.Text, " ")

			if len(lexemes) < 2 {
				continue
			} else if !strings.HasPrefix(lexemes[1], "https://auto.ria.com/search") {
				continue
			}

			autoRia := autoria_api.NewAPI(apiKey)

			params, err := autoria_api.ParseCarSearchParams(lexemes[1])

			if err != nil {
				continue
			}

			search := autoRia.GetSearchCars(params...)

			msg := tgbotapi.NewMessage(msg.Chat.ID, "Hi")
			msg.Text = strings.Join(search.Result.SearchResult.CarsIDs, " ")
			bot.Send(msg)
		}
	}

}
