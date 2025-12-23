package config

import (
	"log"

	"github.com/critma/prodfiles/internal/store"
	"github.com/spf13/viper"
)

type Application struct {
	Config *Config
	Logger *log.Logger
	Store  store.Storage
}

type Config struct {
	Addr     string `mapstructure:"ADDRESS"`
	Port     string `mapstructure:"PORT"`
	BasePath string `mapstructure:"LOCAL_STORE_BASE_PATH"`

	LogLevel string `mapstructure:"LOG_LEVEL"`
}

func SetConfig() (config *Config, err error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()
	if err = viper.ReadInConfig(); err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
