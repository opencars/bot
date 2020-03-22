package sqlstore_test

import (
	"os"
	"testing"

	"github.com/opencars/bot/pkg/config"
)

var (
	conf *config.Store
)

func TestMain(m *testing.M) {
	conf = &config.Store{
		Host:     os.Getenv("DATABASE_HOST"),
		Port:     5432,
		User:     "postgres",
		Password: os.Getenv("DATABASE_PASSWORD"),
		Database: "bot_test",
	}

	if conf.Host == "" {
		conf.Host = "127.0.0.1"
	}

	code := m.Run()
	os.Exit(code)
}
