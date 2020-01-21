package main

import (
	"flag"
	"log"

	"github.com/opencars/bot/internal/bot"
	"github.com/opencars/bot/internal/subscription"
	"github.com/opencars/bot/pkg/autoria"
	"github.com/opencars/bot/pkg/config"
	"github.com/opencars/bot/pkg/env"
	"github.com/opencars/bot/pkg/handlers"
	"github.com/opencars/bot/pkg/openalpr"
	"github.com/opencars/toolkit"
)

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

	app := bot.New()
	app.HandleFuncRegexp(`^\p{L}{2}\d{4}\p{L}{2}$`, openCarsHandler.PlatesHandler)
	app.HandleFuncRegexp(`^\p{L}{3}\d{6}$`, openCarsHandler.RegistrationHandler)
	app.HandleFuncRegexp(`^/auto_[0-9]+$`, autoRiaHandler.CarInfoHandler)
	app.HandleFunc("/start", handlers.StartHandler)
	app.HandleFunc("/follow", autoRiaHandler.FollowHandler)
	app.HandleFunc("/stop", autoRiaHandler.StopHandler)
	app.HandleFunc("/number", openCarsHandler.PlatesHandler)
	app.HandleFunc("/registration", openCarsHandler.RegistrationHandler)
	app.HandleFunc("/vin", openCarsHandler.ReportByVIN)
	app.HandlePhoto(openCarsHandler.PhotoHandler)

	log.Println("Listening on port", port)
	if err := app.Listen(host, port); err != nil {
		log.Fatal(err)
	}
}
