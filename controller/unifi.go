package controller

import (
	"fmt"
	"github.com/sequoiia/unifi-proper-portal/model"
	"github.com/sequoiia/unifi-proper-portal/unifi"
	"log"
	"net/http"
)

var Uclient *unifi.Client

func UniFiClientInit() {
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
}

func UniFiCallback(w http.ResponseWriter, r *http.Request) {
	UniFiClientInit()

	guestDetails := model.GetUniFiGuestFromCallback(r)

	var cookies []*http.Cookie = []*http.Cookie{
		{
			Name:   "UPP_clientmac",
			Value:  guestDetails.ClientMacAddress,
			Path:   "/",
			Domain: Config.Domain,
		},
		{
			Name:   "UPP_apmac",
			Value:  guestDetails.AccessPointMacAddress,
			Path:   "/",
			Domain: Config.Domain,
		},
		{
			Name:   "UPP_redirecturl",
			Value:  guestDetails.RedirectUrl,
			Path:   "/",
			Domain: Config.Domain,
		},
		{
			Name:   "UPP_ssid",
			Value:  guestDetails.Ssid,
			Path:   "/",
			Domain: Config.Domain,
		},
		{
			Name:   "UPP_timestamp",
			Value:  fmt.Sprintf("%d", guestDetails.Timestamp),
			Path:   "/",
			Domain: Config.Domain,
		},
	}

	for i := 0; i < len(cookies); i++ {
		http.SetCookie(w, cookies[i])
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
