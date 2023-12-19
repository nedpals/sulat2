package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type SiteController struct {
	*chi.Mux
}

func NewSiteController() *SiteController {
	r := &SiteController{
		Mux: chi.NewRouter(),
	}

	r.Get("/", wrapHandler(r.getSites))
	r.Route("/{siteId}", func(sr chi.Router) {
		sr.Use(getSiteCtx)
		sr.Get("/", wrapHandler(r.getSite))
		sr.With(getQueryCtx).Get("/search", wrapHandler(searchWithinSite))
		sr.Mount("/collections", NewCollectionController())
	})

	return r
}

func (c *SiteController) getSites(w http.ResponseWriter, r *http.Request) error {
	inst := getCurrentInstance(r)
	sites, err := inst.Sites()
	if err != nil {
		return err
	}
	return returnJson(w, sites)
}

func (c *SiteController) getSite(w http.ResponseWriter, r *http.Request) error {
	site := getCurrentSite(r)
	return returnJson(w, site)
}
