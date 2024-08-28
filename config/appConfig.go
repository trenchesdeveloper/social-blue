package config

import (
	"github.com/spf13/viper"
)

type AppConfig struct {
	ServerPort   string `mapstructure:"HTTP_PORT"`
	DSN          string `mapstructure:"DSN"`
	MigrationURL string `mapstructure:"MIGRATION_URL"`
	DBSource     string `mapstructure:"DB_SOURCE"`
	AppSecret    string `mapsctructure:"APP_SECRET"`
	DBdriver     string `mapstructure:"DB_DRIVER"`
}

func LoadConfig(path string) (*AppConfig, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err := viper.ReadInConfig()

	if err != nil {
		return nil, err
	}

	var config AppConfig

	err = viper.Unmarshal(&config)

	if err != nil {
		return nil, err
	}

	return &config, nil

}
