package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/opencars/toolkit"

	"github.com/opencars/bot/pkg/bot"
	"github.com/opencars/bot/pkg/config"
	"github.com/opencars/bot/pkg/domain/operation"
	"github.com/opencars/bot/pkg/domain/registration"
	"github.com/opencars/bot/pkg/domain/vehicle"
	"github.com/opencars/bot/pkg/env"
	"github.com/opencars/bot/pkg/logger"
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

	// store, err := sqlstore.New(&conf.Store)
	// if err != nil {
	// 	logger.Errorf("sql: %s", err)
	// }

	openCarsURL := env.MustFetch("OPEN_CARS_URL")
	authToken := env.MustFetch("OPEN_CARS_API_KEY")

	telegramToken := env.MustFetch("TELEGRAM_TOKEN")

	client := toolkit.New(openCarsURL, authToken)
	o := operation.NewService(client)
	r := registration.NewService(client)
	v := vehicle.NewService(o, r)

	bot, err := bot.NewBot(v, telegramToken)
	if err != nil {
		logger.Fatalf("%s", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-stop
		cancel()
	}()

	if err := bot.Process(ctx); err != nil {
		logger.Fatalf("%s", err)
	}
}
