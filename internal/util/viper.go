package util

import (
	"image_processing_service/internal/models"

	"github.com/spf13/viper"
)


func LoadConfig(path string) (config models.Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err=viper.ReadInConfig()
	if err!=nil{
		return
	}
	err=viper.Unmarshal(&config)
	return
}
