package util

import (
	"net/http"

	commonmodels "dev.azure.com/technovert-vso/Zappr/_git/Zappr/models"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
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

	var handler http.Handler

	handler = r

	handler = cors.AllowAll().Handler(r)
	
	return handler
}