package main

import (
	"fmt"
	"log"
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
	path := env.Get("DATA_PATH", "/tmp/bot.json")
	port := env.Get("PORT", "8080")
	host := env.MustGet("HOST")

	recognizerURL := env.MustGet("RECOGNIZER_URL")
	openCarsURL := env.MustGet("OPEN_CARS_URL")
	autoRiaURL := env.MustGet("AUTO_RIA_TOKEN")

	tbot := bot.New(path, recognizerURL, openCarsURL)

	tbot.HandleFunc("/start", func(api *tgbotapi.BotAPI, msg *tgbotapi.Message) {
		text := fmt.Sprintf("Привіт, %s!", msg.Chat.FirstName)
		if err := bot.Send(api, msg.Chat, text); err != nil {
			log.Printf("send error: %s", err.Error())
		}
	})

	autoRiaHandler := handlers.AutoRiaHandler{
		API:           autoria.New(autoRiaURL),
		Recognizer:    &openalpr.API{URL: recognizerURL},
		Storage:       &opencars.API{URI: openCarsURL},
		Subscriptions: make(map[int64]*subscription.Subscription),
		FilePath:      path,
	}

	openCarsHandler := handlers.OpenCarsHandler{
		OpenCars: &opencars.API{URI: openCarsURL},
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

	log.Panic(tbot.Listen(host, port))
}
