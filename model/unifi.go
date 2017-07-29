package model

import (
	"log"
	"net/http"
	"strconv"
)

type UniFiCallbackGuest struct {
	ClientMacAddress      string
	AccessPointMacAddress string
	Timestamp             int64
	RedirectUrl           string
	Ssid                  string
}

func GetUniFiGuest(r *http.Request) UniFiCallbackGuest {
	u := UniFiCallbackGuest{}
	u.ClientMacAddress = r.URL.Query().Get("id")
	u.AccessPointMacAddress = r.URL.Query().Get("ap")

	var err error
	u.Timestamp, err = strconv.ParseInt(r.URL.Query().Get("t"), 10, 64)
	if err != nil {
		u.Timestamp = 0
		log.Println(err)
	}

	u.RedirectUrl = r.URL.Query().Get("uri")
	u.Ssid = r.URL.Query().Get("ssid")

	return u
}
