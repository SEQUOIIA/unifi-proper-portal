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

type UniFiLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UniFiGuestAuthoriseRequest struct {
	Mac     string   `json:"mac"`
	Minutes float64  `json:"minutes"`
	Up      *float64 `json:"up,omitempty"`
	Down    *float64 `json:"down,omitempty"`
	Bytes   *float64 `json:"bytes,omitempty"`
}

func GetUniFiGuestFromCallback(r *http.Request) UniFiCallbackGuest {
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
