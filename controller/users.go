package controller

import (
	"net/http"
	"fmt"
)

func UsersView(w http.ResponseWriter, r * http.Request) {
	var listHTML string = "<html>"

	for k, v := range Users {
		listHTML = listHTML + fmt.Sprintf("Name: %s | Email: %s | ID: %s<br>", v.Name, v.Email, k)
	}

	listHTML = listHTML + "</html>"

	w.WriteHeader(200)
	w.Write([]byte(listHTML))
}