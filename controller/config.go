package controller

import (
	"github.com/spf13/viper"
	"log"
)

type ConfigS struct {
	ClientId     string
	ClientSecret string
}

var Config ConfigS

func LoadConfig() {
	viper.SetConfigFile("config.yaml")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal(err)
	}

	Config = ConfigS{}
	Config.ClientId = viper.GetString("clientid")
	Config.ClientSecret = viper.GetString("clientsecret")
}
