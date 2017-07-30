package controller

import (
	"github.com/spf13/viper"
	"log"
)

type ConfigS struct {
	ClientId     string
	ClientSecret string
	UniFi        struct {
		Username     string
		Password     string
		BaseUrl      string
		Site         string
		InsecureMode bool
	}
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
	unifiVars := viper.GetStringMap("unifi")
	Config.UniFi.Username = unifiVars["username"].(string)
	Config.UniFi.Password = unifiVars["password"].(string)
	Config.UniFi.BaseUrl = unifiVars["baseurl"].(string)
	Config.UniFi.Site = unifiVars["site"].(string)
	Config.UniFi.InsecureMode = unifiVars["insecuremode"].(bool)
}
