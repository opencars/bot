package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/shal/opencars-bot/internal/bot"
	"github.com/shal/opencars-bot/internal/subscription"
	"github.com/shal/opencars-bot/pkg/autoria"
	"github.com/shal/opencars-bot/pkg/env"
	"github.com/shal/opencars-bot/pkg/handlers"
	"github.com/shal/opencars-bot/pkg/openalpr"
	"github.com/shal/opencars-bot/pkg/opencars"
)

func StartHandler(msg *bot.Message) {
	if err := msg.SetStatus(bot.ChatTyping); err != nil {
		log.Printf("action error: %s", err.Error())
	}

	text := fmt.Sprintf("Привіт, %s!", msg.Chat().FirstName)
	if err := msg.Send(text); err != nil {
		log.Printf("send error: %s", err.Error())
	}
}
//
//func PhotoHandler(api *tgbotapi.BotAPI, msg *tgbotapi.Message) {
//
//}

// Looks still not very beautiful.
// TODO: Consider refactoring "bot" library to make it's usage much cleaner.
func main() {
	path := env.Fetch("DATA_PATH", "/etc/bot.json")
	port := env.Fetch("PORT", "8080")
	host := env.MustFetch("HOST")

	recognizerURL := env.MustFetch("RECOGNIZER_URL")
	openCarsURL := env.MustFetch("OPEN_CARS_URL")
	autoRiaURL := env.MustFetch("AUTO_RIA_TOKEN")

	tbot := bot.New(path, recognizerURL, openCarsURL)

	tbot.HandleFunc("/start", StartHandler)

	autoRiaHandler := handlers.AutoRiaHandler{
		API:           autoria.New(autoRiaURL),
		Recognizer:    &openalpr.API{URI: recognizerURL},
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

	tbot.HandleFunc(bot.PhotoEvent, func(msg *bot.Message) {
		msg.Send("Nice photo!")
		// TODO: Implement logic of plates detection.
	})

	if err := tbot.Listen(host, port); err != nil {
		log.Panic(err)
	}
}
