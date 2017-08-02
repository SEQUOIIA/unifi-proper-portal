package controller

import (
	"bytes"
	"encoding/json"
	"github.com/sequoiia/unifi-proper-portal/model"
	"log"
	"net/http"
)

func ApiStatus(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Error        int    `json:"error"`
		ErrorMessage string `json:"message"`
		Data         struct {
			Authorisation uint8 `json:"authorisation"`
		} `json:"data"`
	}
	payload.Error = 0

	id, err := model.GetCookieValue(r, "UPP_ID")
	if err != nil {
		payload.Error = 1
		payload.ErrorMessage = "User couldn't be found."
	} else {
		if user, ok := Users[id]; ok {
			payload.Data.Authorisation = user.Authorised
		} else {
			payload.Error = 1
			payload.ErrorMessage = "User couldn't be found."
		}
	}

	buf := bytes.NewBuffer(nil)

	err = json.NewEncoder(buf).Encode(payload)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(200)
	w.Write(buf.Bytes())
}
