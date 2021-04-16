package main

import (
	"context"
	"flag"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/opencars/bot/pkg/bot"
	"github.com/opencars/bot/pkg/config"
	"github.com/opencars/bot/pkg/domain/alpr"
	"github.com/opencars/bot/pkg/domain/vehicle"
	"github.com/opencars/bot/pkg/logger"
	"github.com/opencars/bot/pkg/store/sqlstore"
	"github.com/opencars/toolkit"
)

func main() {
	cfg := flag.String("config", "config/config.yaml", "Path to the configuration file")
	port := flag.Int("port", 8080, "Port of the server")

	flag.Parse()

	conf, err := config.New(*cfg)
	if err != nil {
		logger.Fatalf("config: %v", err)
	}

	logger.NewLogger(logger.LogLevel(conf.Log.Level), conf.Log.Mode == "dev")

	s, err := sqlstore.New(&conf.Database)
	if err != nil {
		logger.Fatalf("store: %v", err)
	}

	// TODO: Implement gRPC API for ALPR service.
	r := alpr.NewService(
		toolkit.New(conf.ALPR.URL, conf.ALPR.Token),
	)

	svc, err := vehicle.NewService(conf.GRPC.Vehicle.Address(), r)
	if err != nil {
		logger.Fatalf("store: %v", err)
	}

	addr := "0.0.0.0:" + strconv.Itoa(*port)

	b, err := bot.NewBot(&conf.Bot, svc, s.Message(), addr)
	if err != nil {
		logger.Fatalf("declare bot: %s", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := b.Process(ctx); err != nil {
		logger.Fatalf("process bot: %s", err)
	}
}
