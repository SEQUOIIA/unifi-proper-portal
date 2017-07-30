package controller

import (
	"encoding/json"
	"fmt"
	"github.com/sequoiia/unifi-proper-portal/model"
	"log"
	"net/http"
)

var graphBaseURL string = "https://graph.facebook.com"

var HttpClient *http.Client = http.DefaultClient

func convertCodeToShortLivedToken(code string) *model.OAuth_accessTokenResponse {
	req, err := http.NewRequest("GET", graphBaseURL+"/v2.10/oauth/access_token", nil)
	if err != nil {
		log.Fatal(err)
	}

	queries := req.URL.Query()
	queries.Add("client_id", Config.ClientId)
	queries.Add("redirect_uri", fmt.Sprintf("http://%s/social/fb/auth", Config.Domain))
	queries.Add("client_secret", Config.ClientSecret)
	queries.Add("code", code)
	req.URL.RawQuery = queries.Encode()

	resp, err := HttpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	jDecoder := json.NewDecoder(resp.Body)
	var accessTokenResponse *model.OAuth_accessTokenResponse = &model.OAuth_accessTokenResponse{}
	err = jDecoder.Decode(accessTokenResponse)
	if err != nil {
		log.Fatal(err)
	}

	return accessTokenResponse
}

func convertShortLivedTokenToLongLivedToken(atr *model.OAuth_accessTokenResponse) {
	req, err := http.NewRequest("GET", graphBaseURL+"/oauth/access_token", nil)
	if err != nil {
		log.Fatal(err)
	}

	queries := req.URL.Query()
	queries.Add("client_id", Config.ClientId)
	queries.Add("grant_type", "fb_exchange_token")
	queries.Add("client_secret", Config.ClientSecret)
	queries.Add("fb_exchange_token", atr.AccessToken)
	req.URL.RawQuery = queries.Encode()

	resp, err := HttpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	jDecoder := json.NewDecoder(resp.Body)
	*atr = model.OAuth_accessTokenResponse{}
	err = jDecoder.Decode(atr)
	if err != nil {
		log.Fatal(err)
	}

	//return atr
}

func getProfile(atr *model.OAuth_accessTokenResponse) *model.FacebookUserNormal {
	req, err := http.NewRequest("GET", graphBaseURL+"/v2.9/me", nil)
	if err != nil {
		log.Fatal(err)
	}

	queries := req.URL.Query()
	queries.Add("access_token", atr.AccessToken)
	queries.Add("fields", "email, id, gender, name, about, work, link, picture")
	req.URL.RawQuery = queries.Encode()

	resp, err := HttpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	jDecoder := json.NewDecoder(resp.Body)
	var profile *model.FacebookUserNormal = &model.FacebookUserNormal{}
	err = jDecoder.Decode(profile)
	if err != nil {
		log.Fatal()
	}

	return profile
}

func newUserCookie(w http.ResponseWriter, id string) {
	/*
		cJar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
		if err != nil {
			log.Fatal(err)
		}

		var cookies []*http.Cookie
	*/
	var idCookie *http.Cookie = &http.Cookie{
		Name:   "UPP_ID",
		Value:  id,
		Path:   "/",
		Domain: Config.Domain,
	}

	http.SetCookie(w, idCookie)
}

func OAuthRedirect(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	unifiDetails := model.GetUniFiGuestCookies(r)

	// Convert code to short-lived access token
	var accessTokenResponse *model.OAuth_accessTokenResponse = convertCodeToShortLivedToken(code)
	log.Println(accessTokenResponse)

	// Exchange the short-lived access token for a long-lived access token
	convertShortLivedTokenToLongLivedToken(accessTokenResponse)

	log.Println(accessTokenResponse)

	fbProfile := getProfile(accessTokenResponse)
	user := model.Client{Name: fbProfile.Name, Email: fbProfile.Email, Authorised: 0, Tokens: model.Tokens{Facebook: accessTokenResponse}}
	user.Id = model.GenerateClientId(user)
	user.Device = unifiDetails.ClientMacAddress

	Users[user.Id] = user
	newUserCookie(w, user.Id)

	//w.WriteHeader(200)
	//w.Write([]byte("gg"))
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
