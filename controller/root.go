package controller

import (
	"net/http"
	"fmt"
	"log"
)

func Root(w http.ResponseWriter, r * http.Request) {
	var output string = ""
	var idFound bool = true
	id, err := r.Cookie("UNIFI_ID")
	if err != nil {
		if err != http.ErrNoCookie {
			log.Fatal(err)
		} else {
			idFound = false
		}
	}

	if idFound {
		if user, ok := Users[id.Value]; ok {
			output = fmt.Sprintf("Name: %s | Email: %s | ID: %s", user.Name, user.Email, user.Id)
		} else {
			idFound = false
		}
	}

	w.WriteHeader(200)
	if !idFound {
		output = fmt.Sprintf("<html>Captive portal page<br><a href=\"https://www.facebook.com/v2.9/dialog/oauth?client_id=%s&redirect_uri=http://localhost:8080/social/fb/auth&response_type=code&scope=public_profile,email,user_friends\">FB Login</a></html>", Config.ClientId)
	}
	w.Write([]byte(output))
}