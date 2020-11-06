package main

import (
	"flag"
	"fmt"

	"github.com/opencars/toolkit"

	"github.com/opencars/bot/internal/bot"
	"github.com/opencars/bot/internal/subscription"
	"github.com/opencars/bot/pkg/autoria"
	"github.com/opencars/bot/pkg/config"
	"github.com/opencars/bot/pkg/env"
	"github.com/opencars/bot/pkg/handlers"
	"github.com/opencars/bot/pkg/logger"
	"github.com/opencars/bot/pkg/store/sqlstore"
)

func main() {
	var configPath string
	var port int

	flag.StringVar(&configPath, "config", "config/config.yaml", "Path to the application configuration file")
	flag.IntVar(&port, "port", 8080, "")

	flag.Parse()

	conf, err := config.New(configPath)
	if err != nil {
		logger.Errorf("config: %s", err)
	}

	store, err := sqlstore.New(&conf.Store)
	if err != nil {
		logger.Errorf("sql: %s", err)
	}

	host := env.MustFetch("HOST")

	openCarsURL := env.MustFetch("OPEN_CARS_URL")
	authToken := env.MustFetch("OPEN_CARS_API_KEY")

	autoRiaHandler := handlers.AutoRiaHandler{
		API:           autoria.New(conf.AutoRia.ApiKey),
		Period:        conf.AutoRia.Period.Duration,
		Toolkit:       toolkit.New(openCarsURL, authToken),
		Subscriptions: make(map[int64]*subscription.Subscription),
	}

	openCarsHandler := handlers.NewOpenCarsHandler(
		toolkit.New(openCarsURL, authToken),
	)

	telegramToken := env.MustFetch("TELEGRAM_TOKEN")

	app, err := bot.New(telegramToken, store)
	if err != nil {
		logger.Fatalf("bot: %s", err)
	}

	app.HandleFuncRegexp(`^\p{L}{2}\d{4}\p{L}{2}$`, openCarsHandler.PlatesHandler)
	app.HandleFuncRegexp(`^\p{L}{3}\d{6}$`, openCarsHandler.RegistrationHandler)
	app.HandleFuncRegexp(`^/auto_[0-9]{8}$`, autoRiaHandler.CarInfoHandler)
	app.HandleFuncRegexp(`^https://auto.ria.com(/uk)?/auto_(.*)_([0-9]{8}).html$`, autoRiaHandler.CarInfoHandler)
	app.HandleFuncRegexp(`^https://auto.ria.com(/uk)?/search/(.*)$`, autoRiaHandler.FollowHandler)
	app.HandleFuncRegexp(`^[A-HJ-NPR-Z0-9]{17}$`, openCarsHandler.ReportByVIN)
	app.HandleFunc("/start", handlers.StartHandler)
	app.HandleFunc("/stop", autoRiaHandler.StopHandler)
	app.HandleFunc("/number", openCarsHandler.PlatesHandler)
	app.HandleFunc("/vin", openCarsHandler.ReportByVIN)

	app.HandlePhoto(openCarsHandler.PhotoHandler)

	addr := fmt.Sprintf(":%d", port)
	logger.Infof("Listening on %s...", addr)

	if err := app.Listen(host, addr); err != nil {
		logger.Fatalf("listen: %s", err)
	}
}
