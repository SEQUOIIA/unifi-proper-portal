package controller

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sequoiia/unifi-proper-portal/model"
	"net/http"
)

var Users map[string]model.Client = make(map[string]model.Client)

func UsersView(w http.ResponseWriter, r *http.Request) {
	var listHTML string = "<html>"

	for k, v := range Users {
		listHTML = listHTML + fmt.Sprintf("Name: %s | Email: %s | ID: %s | Mac: %s", v.Name, v.Email, k, v.Device)

		switch v.Authorised {
		case 0:
			listHTML = listHTML + fmt.Sprintf(" | Authorisation: <a href=\"/users/auth/%s?authorisation=2\">Approve</a> <a href=\"/users/auth/%s?authorisation=1\">Decline</a> <a href=\"/users/delete/%s\">Delete</a><br>", v.Id, v.Id, v.Id)
		case 1:
			listHTML = listHTML + fmt.Sprintf(" | Authorisation: Unauthorised, <a href=\"/users/auth/%s?authorisation=2\">Authorise</a> <a href=\"/users/delete/%s\">Delete</a><br>", v.Id, v.Id)
		case 2:
			listHTML = listHTML + fmt.Sprintf(" | Authorisation: Authorised, <a href=\"/users/auth/%s?authorisation=1\">Unauthorise</a> <a href=\"/users/delete/%s\">Delete</a><br>", v.Id, v.Id)
		}
	}

	listHTML = listHTML + "</html>"

	w.WriteHeader(200)
	w.Write([]byte(listHTML))
}

func UsersDeleteApi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var success bool = false

	if clientid, ok := vars["clientid"]; ok {
		if _, ok := Users[clientid]; ok {
			delete(Users, clientid)
			success = true
		}
	}

	if success {
		w.WriteHeader(200)
		w.Write([]byte("Success"))
	} else {
		w.WriteHeader(400)
		w.Write([]byte("Invalid"))
	}
}

func UsersAuthorisationApi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clientId := vars["clientid"]
	authCode := r.URL.Query().Get("authorisation")
	var set bool = false

	if authCode != "" {
		if user, ok := Users[clientId]; ok {
			switch authCode {
			case "1":
				user.Authorised = 1
				Users[clientId] = user
				set = true
			case "2":
				user.Authorised = 2
				Users[clientId] = user
				set = true

				Uclient.AuthoriseGuest(model.UniFiGuestAuthoriseRequest{
					Mac:     user.Device,
					Minutes: 0,
				})
			default:
				w.WriteHeader(400)
				w.Write([]byte("Invalid"))
			}

		} else {
			w.WriteHeader(400)
			w.Write([]byte("Invalid"))
		}

	} else {
		w.WriteHeader(400)
		w.Write([]byte("Invalid"))
	}

	if set {
		w.WriteHeader(200)
		w.Write([]byte("Success."))
	}
}
