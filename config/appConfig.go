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
	SendgridAPIKey string `mapstructure:"SENDGRID_API_KEY"`
	SendgridFromEmail string `mapstructure:"SENDGRID_FROM_EMAIL"`
	FrontendURL  string `mapstructure:"FRONTEND_URL"`
	MAILTRAP_API_KEY string `mapstructure:"MAILTRAP_API_KEY"`
	BASIC_AUTH_USERNAME string `mapstructure:"BASIC_AUTH_USERNAME"`
	BASIC_AUTH_PASSWORD string `mapstructure:"BASIC_AUTH_PASSWORD"`
}

type MailConfig struct {
	EXP time.Duration
}

func LoadConfig(path string) (*AppConfig, error) {
	// Always load environment variables from the environment
	viper.AutomaticEnv()

	// bind environment variables
	viper.BindEnv("HTTP_PORT", "HTTP_PORT")
	viper.BindEnv("DSN", "DSN")
	viper.BindEnv("MIGRATION_URL", "MIGRATION_URL")
	viper.BindEnv("DB_SOURCE", "DB_SOURCE")
	viper.BindEnv("APP_SECRET, APP_SECRET")
	viper.BindEnv("DB_DRIVER", "DB_DRIVER")
	viper.BindEnv("ENVIRONMENT", "ENVIRONMENT")
	viper.BindEnv("API_URL", "API_URL")
	viper.BindEnv("SENDGRID_API_KEY", "SENDGRID_API_KEY")
	viper.BindEnv("SENDGRID_FROM_EMAIL", "SENDGRID_FROM_EMAIL")
	viper.BindEnv("FRONTEND_URL", "FRONTEND_URL")
	viper.BindEnv("MAILTRAP_API_KEY", "MAILTRAP_API_KEY")
	viper.BindEnv("BASIC_AUTH_USERNAME", "BASIC_AUTH_USERNAME")
	viper.BindEnv("BASIC_AUTH_PASSWORD", "BASIC_AUTH_PASSWORD")

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



	var config AppConfig
	err := viper.Unmarshal(&config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}
