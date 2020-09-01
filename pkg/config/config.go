package config

import (
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	AutoRia *AutoRia `toml:"autoria"`
	Store   *Store   `toml:"store"`
}

type Duration struct {
	time.Duration
}

type AutoRia struct {
	Period Duration `toml:"period"`
	ApiKey string   `toml:"api_key"`
}

type Store struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Database string `toml:"database"`
	SSLMode  string `toml:"ssl_mode"`
}

// UnmarshalText implements encoding.TextUnmarshaler{} for Duration type.
func (d *Duration) UnmarshalText(text []byte) error {
	var err error

	d.Duration, err = time.ParseDuration(string(text))
	return err
}

// New reads application configuration from specified file path.
func New(path string) (*Config, error) {
	var config Config

	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
