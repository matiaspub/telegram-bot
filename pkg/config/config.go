package config

import (
	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
	"os"
)

type Config struct {
	TelegramToken     string `env:"TELEGRAM_TOKEN"`
	TelegramBotUrl    string `mapenv:"TELEGRAM_BOT_URL"`
	PocketConsumerKey string `env:"POCKET_CONSUMER_KEY"`
	AuthServerUrl     string `mapenv:"AUTH_SERVER_URL"`
	DbFile            string `mapstructure:"db_file"`

	Messages Messages
}

type Messages struct {
	Errors
	Responses
}

type Responses struct {
	Start             string `mapstructure:"start"`
	AlreadyAuthorized string `mapstructure:"already_authorized"`
	SavedSuccessfully string `mapstructure:"saved_successfully"`
	UnknownCommand    string `mapstructure:"unknown_command"`
}

type Errors struct {
	InvalidUrl   string `mapstructure:"invalid_url"`
	Unauthorized string `mapstructure:"unauthorized"`
	UnableToSave string `mapstructure:"unable_to_save"`
}

func Init() (*Config, error) {
	viper.AddConfigPath("config")
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("messages.responses", &cfg.Messages.Responses); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("messages.errors", &cfg.Messages.Errors); err != nil {
		return nil, err
	}

	if err := parseEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func parseEnv(cfg *Config) error {
	if err := gotenv.Load(); err != nil {
		return err
	}

	cfg.TelegramToken = os.Getenv("TELEGRAM_TOKEN")
	cfg.TelegramBotUrl = os.Getenv("TELEGRAM_BOT_URL")
	cfg.PocketConsumerKey = os.Getenv("POCKET_CONSUMER_KEY")
	cfg.AuthServerUrl = os.Getenv("AUTH_SERVER_URL")

	return nil
}
