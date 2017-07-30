package controller

import (
	"bytes"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"log"
	"net"
	"net/http"
)

func NewUsersRouter() *mux.Router {
	router := mux.NewRouter()
	users := router.PathPrefix("/users").Subrouter()

	users.Handle("/", negroni.New(
		negroni.Wrap(http.HandlerFunc(UsersView)),
	))

	users.Handle("/auth/{clientid}", negroni.New(
		negroni.Wrap(http.HandlerFunc(UsersAuthorisationApi)),
	))

	users.Handle("/delete/{clientid}", negroni.New(
		negroni.Wrap(http.HandlerFunc(UsersDeleteApi)),
	))
	return router
}

func SubnetFenceMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ip := parseIp(r.RemoteAddr)
	var approved bool = false

	for i := 0; i < len(Config.AllowedSubnets); i++ {
		subnet := Config.AllowedSubnets[i]
		// IPv4 check
		if p4 := ip.To4(); len(p4) == net.IPv4len {
			if subnetp4 := subnet.Ip.To4(); len(subnetp4) == net.IPv4len {
				var valid bool = true
				for i := uint8(0); i <= subnet.Count; i++ {
					if uint8(p4[i]) != uint8(subnetp4[i]) {
						valid = false
					}
				}

				if valid {
					log.Printf("%s accessing %s\n", ip.String(), r.URL.Path)
					approved = true
					next(rw, r)
					break
				}

			}
		}

	}

	if !approved {
		log.Printf("%s denied access to %s\n", ip.String(), r.URL.Path)
		rw.Write([]byte("Access denied."))
		return
	}
}

func parseIp(ipAndPort string) net.IP {
	var rawIpAndPort string = ipAndPort
	var rawIp string = ""

	if rawIpAndPort[0] != '[' { // It's an IPv4 address
		buf := bytes.NewBuffer(nil)
		var portReached bool = false
		var count int = 0

		for !portReached {
			if rawIpAndPort[count] != ':' {
				buf.WriteByte(rawIpAndPort[count])
				count++
			} else {
				portReached = true
			}
		}

		rawIp = buf.String()
	} else { // It's an IPv6 address
		buf := bytes.NewBuffer(nil)
		var portReached bool = false
		var count int = 0

		for !portReached {
			switch rawIpAndPort[count] {
			case ']':
				portReached = true
			case '[':
				count++
			default:
				buf.WriteByte(rawIpAndPort[count])
				count++
			}
		}

		rawIp = buf.String()
	}

	return net.ParseIP(rawIp)
}
