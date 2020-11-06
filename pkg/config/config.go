package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

// Config represents mix of settings for the app.
type Config struct {
	AutoRia AutoRia `yaml:"autoria"`
	Store   Store   `yaml:"database"`
}

// Server represents settings for creating http server.
type Server struct {
	ShutdownTimeout Duration `yaml:"shutdown_timeout"`
	ReadTimeout     Duration `yaml:"read_timeout"`
	WriteTimeout    Duration `yaml:"write_timeout"`
	IdleTimeout     Duration `yaml:"idle_timeout"`
}

// Log represents settings for application logger.
type Log struct {
	Level string `yaml:"level"`
	Mode  string `yaml:"mode"`
}

// AutoRia represents configuration for the AutoRia API.
type AutoRia struct {
	Period Duration `yaml:"period"`
	ApiKey string   `yaml:"api_key"`
}

// Store represents configuration for the storage.
type Store struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	SSLMode  string `yaml:"ssl_mode"`
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
