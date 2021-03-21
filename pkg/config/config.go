package config

import (
	"os"
	"strconv"

	"gopkg.in/yaml.v2"
)

// Config represents mix of settings for the app.
type Config struct {
	Log      Log      `yaml:"log"`
	Server   Server   `yaml:"server"`
	Database Database `yaml:"database"`
	GRPC     GRPC     `yaml:"grpc"`
	Bot      Bot      `yaml:"bot"`
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
	APIKey string   `yaml:"api_key"`
}

// Store represents configuration for the storage.
type Database struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	SSLMode  string `yaml:"ssl_mode"`
}

type GRPC struct {
	Vehicle ServiceGRPC `yaml:"vehicle"`
}

type ServiceGRPC struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Bot struct {
	URL            string `yaml:"url"`
	Token          string `yaml:"token"`
	MaxConnections *int   `yaml:"max_connections"`
}

func (s *ServiceGRPC) Address() string {
	return s.Host + ":" + strconv.Itoa(s.Port)
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
