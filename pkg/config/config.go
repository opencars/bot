package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	AutoRia *AutoRia `yaml:"autoria"`
	Store   *Store   `yaml:"database"`
}

type Duration struct {
	time.Duration
}

type AutoRia struct {
	Period Duration `yaml:"period"`
	ApiKey string   `yaml:"api_key"`
}

type Store struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	SSLMode  string `yaml:"ssl_mode"`
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

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.NewDecoder(f).Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
