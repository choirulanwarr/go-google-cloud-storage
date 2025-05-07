package config

import (
	"fmt"
	"github.com/spf13/viper"
)

func NewViper() *viper.Viper {
	config := viper.New()

	config.AutomaticEnv()
	config.AddConfigPath("./")
	config.SetConfigName(".env")
	config.SetConfigType("env")

	if err := config.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error load env file: %w \n", err))
	}

	return config
}
