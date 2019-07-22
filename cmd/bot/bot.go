package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/opencars/toolkit/sdk"

	"github.com/opencars/bot/internal/bot"
	"github.com/opencars/bot/internal/subscription"
	"github.com/opencars/bot/pkg/autoria"
	"github.com/opencars/bot/pkg/env"
	"github.com/opencars/bot/pkg/handlers"
	"github.com/opencars/bot/pkg/openalpr"
)

func StartHandler(msg *bot.Event) {
	if err := msg.SetStatus(bot.ChatTyping); err != nil {
		log.Printf("action error: %s", err.Error())
	}

	text := fmt.Sprintf("Привіт, %s!", msg.Message.Chat.FirstName)
	if err := msg.Send(text); err != nil {
		log.Printf("send error: %s", err.Error())
	}
}

// Looks still not very beautiful.
// TODO: Consider refactoring "bot" library to make it's usage much cleaner.
func main() {
	path := env.Fetch("DATA_PATH", "/etc/bot.json")
	port := env.Fetch("PORT", "8080")
	host := env.MustFetch("HOST")

	recognizerURL := env.MustFetch("RECOGNIZER_URL")
	openCarsURL := env.MustFetch("OPEN_CARS_URL")
	autoRiaToken := env.MustFetch("AUTO_RIA_TOKEN")

	app := bot.New(path, recognizerURL, openCarsURL)

	app.HandleFunc("/start", StartHandler)

	autoRiaHandler := handlers.AutoRiaHandler{
		API:           autoria.New(autoRiaToken),
		Recognizer:    &openalpr.API{URI: recognizerURL},
		Storage:       sdk.New(openCarsURL),
		Subscriptions: make(map[int64]*subscription.Subscription),
		FilePath:      path,
	}

	openCarsHandler := handlers.OpenCarsHandler{
		OpenCars:   sdk.New(openCarsURL),
		Recognizer: &openalpr.API{URI: recognizerURL},
	}

	app.HandleFunc("/follow", autoRiaHandler.FollowHandler)
	app.HandleFunc("/stop", autoRiaHandler.StopHandler)

	expr, err := regexp.Compile(`^\p{L}{2}\d{4}\p{L}{2}$`)
	if err != nil {
		log.Panic(err)
	}

	app.HandleFuncRegexp(expr, openCarsHandler.PlatesHandler)

	expr, err = regexp.Compile(`^/auto_[0-9]+$`)
	if err != nil {
		log.Panic(err)
	}

	app.HandleFuncRegexp(expr, autoRiaHandler.CarInfoHandler)
	app.HandleFunc(bot.PhotoEvent, openCarsHandler.PhotoHandler)

	// Handler for "/plates" keyword.
	// Usage: /plates AA1234XX.
	app.HandleFunc("/plates", openCarsHandler.PlatesHandler)

	// Handler for "/vin" keyword.
	// Usage: /vin X0X0XXXXXXX0000X0.
	app.HandleFunc("/vin", openCarsHandler.NotImplemented)

	if err := app.Listen(host, port); err != nil {
		log.Panic(err)
	}
}
