package main

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/shal/robot/internal/bot"
	"github.com/shal/robot/internal/subscription"
	"github.com/shal/robot/pkg/autoria"
	"github.com/shal/robot/pkg/env"
	"github.com/shal/robot/pkg/handlers"
	"github.com/shal/robot/pkg/openalpr"
	"github.com/shal/robot/pkg/opencars"
)

// Looks still not very beautiful.
// TODO: Consider refactoring "bot" library to make it's usage much cleaner.
func main() {
	jsonPath := env.Get("BOT_DATA_PATH", "/tmp/bot.json")
	alprURL := env.Get("ALPR_URL", "http://alpr.opencars.pp.ua")
	apiURL := env.Get("API_URL", "http://api.opencars.pp.ua")
	autoRiaURL := env.MustGet("RIA_API_KEY")

	tbot := bot.New(jsonPath, alprURL, apiURL)

	tbot.HandleFunc("/start", func(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
		text := fmt.Sprintf("Привіт, %s!", msg.Chat.FirstName)
		handlers.Send(bot, msg.Chat, text)
	})

	autoRiaHandler := handlers.AutoRiaHandler{
		API:           autoria.New(autoRiaURL),
		Recognizer:    &openalpr.API{URL: alprURL},
		Storage:       &opencars.API{URI: apiURL},
		Subscriptions: make(map[int64]*subscription.Subscription),
		FilePath:      jsonPath,
	}

	openCarsHandler := handlers.OpenCarsHandler{
		OpenCars: &opencars.API{URI: apiURL},
	}

	tbot.HandleFunc("/follow", autoRiaHandler.FollowHandler)
	tbot.HandleFunc("/stop", autoRiaHandler.StopHandler)

	expr, err := regexp.Compile(`^\p{L}{2}\d{4}\p{L}{2}$`)
	if err != nil {
		log.Panic(err)
	}
	tbot.HandleRegexp(expr, openCarsHandler)

	expr, err = regexp.Compile(`^/auto_\d+$`)
	if err != nil {
		log.Panic(err)
	}
	tbot.HandleFuncRegexp(expr, autoRiaHandler.CarInfoHandler)

	// Host is required.
	host, exists := os.LookupEnv("HOST")
	if !exists {
		log.Panic("host is not specified")
	}

	log.Panic(tbot.Listen(host, env.Get("PORT", "8080")))
}
