package main

import (
	"github.com/GeertJohan/go.rice"
	"github.com/gorilla/mux"
	"github.com/sequoiia/unifi-proper-portal/controller"
	"github.com/urfave/negroni"
	"net/http"
	"strings"
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

	router.Handle("/api/status", negroni.New(
		negroni.Wrap(http.HandlerFunc(controller.ApiStatus)),
	))

	router.Handle("/api/voucher", negroni.New(
		negroni.Wrap(http.HandlerFunc(controller.ApiVoucherConsume)),
	))

	// fb
	router.Handle("/social/fb/auth", negroni.New(
		negroni.Wrap(http.HandlerFunc(controller.OAuthRedirect)),
	))

	// unifi callback
	router.Handle("/guest/s/{site}/", negroni.New(
		negroni.Wrap(http.HandlerFunc(controller.UniFiCallback)),
	))

	// static assets
	assets := http.StripPrefix("/assets/", http.FileServer(rice.MustFindBox("static/build").HTTPBox()))
	//router.PathPrefix("/assets/").Handler(assets)
	router.PathPrefix("/assets/").Handler(negroni.New(
		negroni.HandlerFunc(fixContentType),
		negroni.Wrap(assets),
	))

	n := negroni.New(negroni.NewRecovery())
	n.UseHandler(router)
	n.Run("0.0.0.0:8080")
}

func fixContentType(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	path := r.URL.Path[1:]

	if strings.HasSuffix(path, ".js") {
		rw.Header().Set("Content-Type", "application/javascript")
	}

	next(rw, r)
}
