package config

import (
	"github.com/spf13/viper"
	"time"
)

type AppConfig struct {
	ServerPort   string `mapstructure:"HTTP_PORT"`
	DSN          string `mapstructure:"DSN"`
	MigrationURL string `mapstructure:"MIGRATION_URL"`
	DBSource     string `mapstructure:"DB_SOURCE"`
	AppSecret    string `mapsctructure:"APP_SECRET"`
	DBdriver     string `mapstructure:"DB_DRIVER"`
	Environment  string `mapstructure:"ENVIRONMENT"`
	ApiUrl       string `mapstructure:"API_URL"`
}

type MailConfig struct {
	EXP time.Duration
}

func LoadConfig(path string) (*AppConfig, error) {
	// Check if environment is set to production
	if viper.GetString("ENVIRONMENT") != "production" {
		viper.AddConfigPath(path)
		viper.SetConfigName("app")
		viper.SetConfigType("env")

		err := viper.ReadInConfig()
		if err != nil {
			return nil, err
		}
	}

	// Always load environment variables from the environment
	viper.AutomaticEnv()

	var config AppConfig
	err := viper.Unmarshal(&config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}
