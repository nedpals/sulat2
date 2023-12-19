package rest_server

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/nedpals/sulatcms/sulat"
)

func NewRestRouter(rootInst *sulat.Instance) *chi.Mux {
	r := chi.NewRouter()

	r.Use(
		middleware.Logger,
		middleware.Recoverer,
		cors.Handler(cors.Options{
			AllowedOrigins: []string{"*"},
		}),
		getInstanceCtx(rootInst),
	)

	r.Route("/api", func(r chi.Router) {
		r.Mount("/sites", NewSiteController())
		r.Mount("/collections", NewCollectionController())
		r.Mount("/data-sources", NewDataSourceController())
	})

	return r
}

func StartServer(rootInst *sulat.Instance, defaultPort string) error {
	if len(defaultPort) == 0 {
		defaultPort = "3000"
	}

	port, portExists := os.LookupEnv("PORT")
	if !portExists {
		port = defaultPort
	}

	r := NewRestRouter(rootInst)
	log.Printf("Starting server on port %s\n", port)
	return http.ListenAndServe(":"+port, r)
}
