package controller

import (
	"fmt"
	"github.com/sequoiia/unifi-proper-portal/model"
	"log"
	"net/http"
)

func Root(w http.ResponseWriter, r *http.Request) {
	var output string = ""
	var idFound bool = true
	id, err := r.Cookie("UPP_ID")
	if err != nil {
		if err != http.ErrNoCookie {
			log.Fatal(err)
		} else {
			idFound = false
		}
	}

	if idFound {
		if user, ok := Users[id.Value]; ok {
			switch user.Authorised {
			case 2:
				fbProfile := getProfile(user.Tokens.Facebook)
				output = fmt.Sprintf("<html><img src=\"%s\">Hi %s. You're authorised to use this network. Proceed as you see fit.</html>", fbProfile.Picture.Data.Url, user.Name)
			case 1:
				output = fmt.Sprintf("<html>Hi %s. You're unauthorised to use this network. Contact the administrator of this network for further information.</html>", user.Name)
			case 0:
				output = "<html>Awaiting authorisation.</html>"
			}

			//output = fmt.Sprintf("Name: %s | Email: %s | ID: %s", user.Name, user.Email, user.Id)
		} else {
			idFound = false
		}
	}

	w.WriteHeader(200)
	if !idFound {
		unifiDetails := model.GetUniFiGuestCookies(r)
		log.Println(unifiDetails.ClientMacAddress)
		output = fmt.Sprintf("<html>Captive portal page<br><a href=\"https://www.facebook.com/v2.9/dialog/oauth?client_id=%s&redirect_uri=http://%s/social/fb/auth&response_type=code&scope=public_profile,email,user_friends\">FB Login</a></html>", Config.ClientId, Config.Domain)
	}
	w.Write([]byte(output))
}
