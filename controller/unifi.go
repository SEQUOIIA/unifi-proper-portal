package controller

import (
	"github.com/sequoiia/unifi-proper-portal/model"
	"github.com/sequoiia/unifi-proper-portal/unifi"
	"log"
	"net/http"
)

var Uclient *unifi.Client

func UniFiCallback(w http.ResponseWriter, r *http.Request) {
	if Uclient == nil {
		Uclient = unifi.NewClient(&unifi.ClientConfig{
			Username:        Config.UniFi.Username,
			Password:        Config.UniFi.Password,
			UniFiBaseUrlRaw: Config.UniFi.BaseUrl,
			Site:            Config.UniFi.Site,
			InsecureMode:    Config.UniFi.InsecureMode,
		})

		err := Uclient.Login()
		if err != nil {
			log.Println(err)
		}
	}

	guestDetails := model.GetUniFiGuestFromCallback(r)
}
