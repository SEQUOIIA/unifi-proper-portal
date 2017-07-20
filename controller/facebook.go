package controller

import (
	"net/http"
	"log"
	"encoding/json"
	"faker.clausen/root/unifi-proper-portal/model"
)

var graphBaseURL string = "https://graph.facebook.com"

var HttpClient * http.Client = http.DefaultClient

var Users map[string]model.Client = make(map[string]model.Client)

type oAuth_accessTokenResponse struct {
	AccessToken		string	`json:"access_token"`
	TokenType		string	`json:"token_type"`
	ExpiresIn		int64	`json:"expires_in"`
}

func convertCodeToShortLivedToken(code string) * oAuth_accessTokenResponse {
	req, err := http.NewRequest("GET", graphBaseURL + "/v2.10/oauth/access_token", nil)
	if err != nil {
		log.Fatal(err)
	}

	queries := req.URL.Query()
	queries.Add("client_id", Config.ClientId)
	queries.Add("redirect_uri", "http://localhost:8080/social/fb/auth")
	queries.Add("client_secret", Config.ClientSecret)
	queries.Add("code", code)
	req.URL.RawQuery = queries.Encode()

	resp, err := HttpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}


	jDecoder := json.NewDecoder(resp.Body)
	var accessTokenResponse * oAuth_accessTokenResponse = &oAuth_accessTokenResponse{}
	err = jDecoder.Decode(accessTokenResponse)
	if err != nil {
		log.Fatal(err)
	}

	return accessTokenResponse
}

func convertShortLivedTokenToLongLivedToken(atr * oAuth_accessTokenResponse) {
	req, err := http.NewRequest("GET", graphBaseURL + "/oauth/access_token", nil)
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
	*atr = oAuth_accessTokenResponse{}
	err = jDecoder.Decode(atr)
	if err != nil {
		log.Fatal(err)
	}

	//return atr
}

func getProfile(atr * oAuth_accessTokenResponse) * model.FacebookUserNormal{
	req, err := http.NewRequest("GET", graphBaseURL + "/v2.9/me", nil)
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
	var profile * model.FacebookUserNormal = &model.FacebookUserNormal{}
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
		Name: "UNIFI_ID",
		Value: id,
		Path: "/",
		Domain: "localhost",
	}

	http.SetCookie(w, idCookie)
}

func OAuthRedirect(w http.ResponseWriter, r * http.Request) {
	code := r.URL.Query().Get("code")

	// Convert code to short-lived access token
	var accessTokenResponse * oAuth_accessTokenResponse = convertCodeToShortLivedToken(code)
	log.Println(accessTokenResponse)

	// Exchange the short-lived access token for a long-lived access token
	convertShortLivedTokenToLongLivedToken(accessTokenResponse)

	log.Println(accessTokenResponse)

	fbProfile := getProfile(accessTokenResponse)
	user := model.Client{Name: fbProfile.Name, Email: fbProfile.Email}
	user.Id = model.GenerateClientId(user)


	Users[user.Id] = user
	newUserCookie(w, user.Id)

	w.WriteHeader(200)
	w.Write([]byte("gg"))
}