package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"

	"github.com/opencars/toolkit"

	"github.com/opencars/bot/internal/bot"
	"github.com/opencars/bot/internal/subscription"
	"github.com/opencars/bot/pkg/autoria"
	"github.com/opencars/bot/pkg/config"
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

func main() {
	var configPath string

	flag.StringVar(&configPath, "config", "config/config.toml", "Path to the application configuration file")

	flag.Parse()

	conf, err := config.New(configPath)
	if err != nil {
		log.Fatal(err)
	}

	port := env.Fetch("PORT", "8080")
	host := env.MustFetch("HOST")

	recognizerURL := env.MustFetch("RECOGNIZER_URL")
	openCarsURL := env.MustFetch("OPEN_CARS_URL")
	authToken := env.MustFetch("OPEN_CARS_API_KEY")

	app := bot.New()

	app.HandleFunc("/start", StartHandler)

	autoRiaHandler := handlers.AutoRiaHandler{
		API:           autoria.New(conf.AutoRia.ApiKey),
		Period:        conf.AutoRia.Period.Duration,
		ApiKey:        conf.AutoRia.ApiKey,
		Recognizer:    &openalpr.API{URI: recognizerURL},
		Toolkit:       toolkit.New(openCarsURL, authToken),
		Subscriptions: make(map[int64]*subscription.Subscription),
	}

	openCarsHandler := handlers.NewOpenCarsHandler(
		toolkit.New(openCarsURL, authToken),
		&openalpr.API{URI: recognizerURL},
	)

	expr, err := regexp.Compile(`^\p{L}{2}\d{4}\p{L}{2}$`)
	if err != nil {
		log.Panic(err)
	}
	app.HandleFuncRegexp(expr, openCarsHandler.PlatesHandler)

	expr, err = regexp.Compile(`^\p{L}{3}\d{6}$`)
	if err != nil {
		log.Panic(err)
	}
	app.HandleFuncRegexp(expr, openCarsHandler.RegistrationHandler)

	expr, err = regexp.Compile(`^/auto_[0-9]+$`)
	if err != nil {
		log.Panic(err)
	}
	app.HandleFuncRegexp(expr, autoRiaHandler.CarInfoHandler)

	app.HandleFunc("/follow", autoRiaHandler.FollowHandler)
	app.HandleFunc("/stop", autoRiaHandler.StopHandler)
	app.HandleFunc("/plates", openCarsHandler.PlatesHandler)
	app.HandleFunc("/registration", openCarsHandler.RegistrationHandler)
	app.HandleFunc("/vin", openCarsHandler.NotImplemented)

	app.HandlePhoto(openCarsHandler.PhotoHandler)

	log.Println("Listening on port", port)
	if err := app.Listen(host, port); err != nil {
		log.Panic(err)
	}
}