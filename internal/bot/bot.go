// Copyright (C) 2019 Ali Shanaakh, github@shanaakh.pro
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package bot

import (
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/shal/robot/internal/app"
	"github.com/shal/robot/pkg/env"
)

func Run() {
	jsonPath := env.Get("BOT_DATA_PATH", "/tmp/bot.json")
	alprURL  := env.Get("ALPR_URL", "http://alpr.robot.shanaakh.pro")

	botApp := app.NewApp(jsonPath, alprURL)

	// Bot updates configuration.
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates, err := botApp.Bot.GetUpdatesChan(updateConfig)

	if err != nil {
		log.Println(err.Error())
	}

	// Triggers, when new update pushed to channel.
	for update := range updates {
		botApp.HandleUpdate(update)
	}
}
