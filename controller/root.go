package controller

import (
	"net/http"
	"fmt"
)

func Root(w http.ResponseWriter, r * http.Request) {
	w.WriteHeader(200)
	output := fmt.Sprintf("<html>Captive portal page<br><a href=\"https://www.facebook.com/v2.9/dialog/oauth?client_id=%s&redirect_uri=http://localhost:8080/social/fb/auth&response_type=code&scope=public_profile,email,user_friends\">FB Login</a></html>", Config.ClientId)
	w.Write([]byte(output))
}