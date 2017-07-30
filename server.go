package main

import (
	"github.com/gorilla/mux"
	"github.com/sequoiia/unifi-proper-portal/controller"
	"github.com/urfave/negroni"
	"net/http"
)

func main() {
	controller.LoadConfig() // Load config.yaml
	router := mux.NewRouter()

	router.Handle("/", negroni.New(
		negroni.Wrap(http.HandlerFunc(controller.Root)),
	))

	router.Handle("/users", negroni.New(
		negroni.HandlerFunc(controller.SubnetFenceMiddleware),
		negroni.Wrap(http.HandlerFunc(controller.UsersView)),
	))

	router.PathPrefix("/users").Handler(negroni.New(
		negroni.HandlerFunc(controller.SubnetFenceMiddleware),
		negroni.Wrap(controller.NewUsersRouter()),
	))

	// fb
	router.Handle("/social/fb/auth", negroni.New(
		negroni.Wrap(http.HandlerFunc(controller.OAuthRedirect)),
	))

	// unifi callback
	router.Handle("/guest/s/{site}/", negroni.New(
		negroni.Wrap(http.HandlerFunc(controller.UniFiCallback)),
	))

	n := negroni.New(negroni.NewRecovery())
	n.UseHandler(router)
	n.Run("0.0.0.0:8080")
}
