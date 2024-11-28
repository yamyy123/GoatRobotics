package utils

import (
	"fmt"

	"github.com/spf13/viper"
)

func LoadConfig() error {
	viper.SetConfigName("config") 
	viper.SetConfigType("json")  
	viper.AddConfigPath(".")     

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	return nil
}