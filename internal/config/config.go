package config

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env       string       `yaml:"env" env-default:"local"`
	DbURI     string       `yaml:"db_uri" env-required:"true"`
	RedisADDR string       `yaml:"redis_addr" env-required:"true"`
	GRPC      GRPCConfig   `yaml:"grpc"`
	Clients   ClientConfig `yaml:"clients"`
	AppSecret string       `yaml:"app_secret" env-required:"true" env:"APP_SECRET"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type Client struct {
	Address      string        `yaml:"address"`
	Timeout      time.Duration `yaml:"timeout"`
	RetriesCount int           `yaml:"retriesCount"`
}

type ClientConfig struct {
	SSO Client `yaml:"sso"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("configPath is empty")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exists: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

// fetchConfigPath fetches config path from command line flag or environment variable.
// Priority: flag > env > default.
// Default value is empty string.
func fetchConfigPath() string {
	var configPath string

	// --config="path/to/config.yaml"
	flag.StringVar(&configPath, "config", "", "path to config file")
	flag.Parse()

	if configPath == "" {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("cannot load .env file: %s", err)
		}

		configPath = os.Getenv("CONFIG_PATH")
	}

	return configPath
}
