package controller

import (
	"github.com/sequoiia/unifi-proper-portal/model"
	"html/template"
	"log"
	"net/http"
)

var rootPost0 *template.Template
var rootPost1 *template.Template
var rootPost2 *template.Template
var rootPre *template.Template

func Root(w http.ResponseWriter, r *http.Request) {
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
			var data struct {
				UserName          string
				ProfilePictureUrl string
				Authorised        uint8
				Name              string
				Subtext           string
			}
			data.UserName = user.Name
			data.Authorised = user.Authorised
			fbProfile := getProfile(user.Tokens.Facebook)
			data.ProfilePictureUrl = fbProfile.Picture.Data.Url
			data.Name = Config.Custom.Name
			data.Subtext = Config.Custom.Subtext

			switch user.Authorised {
			case 2:
				templateInit(&rootPost2, "rootPost", STATICPATH+"views/rootPost_2.html")
				rootPost2.Execute(w, data)
			case 1:
				templateInit(&rootPost1, "rootPost", STATICPATH+"views/rootPost_1.html")
				rootPost1.Execute(w, data)
			case 0:
				templateInit(&rootPost0, "rootPost", STATICPATH+"views/rootPost_0.html")
				rootPost0.Execute(w, data)
			}

		} else {
			idFound = false
		}
	}

	//w.WriteHeader(200)
	if !idFound {
		templateInit(&rootPre, "rootPre", STATICPATH+"views/rootPre.html")

		unifiDetails, err := model.GetUniFiGuestCookies(r)
		if err != nil {
			if err != http.ErrNoCookie {
				log.Fatal(err)
			}
		} else {
			log.Println(unifiDetails.ClientMacAddress)
		}

		var data struct {
			ClientId    string
			RedirectUri string
			Name        string
			Subtext     string
		}
		data.ClientId = Config.ClientId
		data.RedirectUri = Config.Domain
		data.Name = Config.Custom.Name
		data.Subtext = Config.Custom.Subtext
		err = rootPre.Execute(w, data)
		if err != nil {
			log.Fatal(err)
		}
	}
	//w.Write([]byte(output))
}

func templateInit(t **template.Template, name string, path string) {
	if *t == nil {
		var err error
		*t, err = template.ParseFiles(path)
		if err != nil {
			log.Fatal(err)
		}
	}
}
