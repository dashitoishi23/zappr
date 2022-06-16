package util

import (
	"net/http"

	commonmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/cmd/models"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func RootHttpHandler(servers []commonmodels.HttpServerConfig) http.Handler {
	r := mux.NewRouter()

	for _, server := range servers {
		if server.NeedsAuth {
			r.Handle(server.Route, negroni.New(
				negroni.HandlerFunc(AuthHandler),
				negroni.Wrap(server.Server),
			)).Methods(server.Methods...)
		} else {
			r.Handle(server.Route, server.Server).Methods(server.Methods...)
		}
	}

	return r
}