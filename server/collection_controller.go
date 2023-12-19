package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nedpals/sulatcms/sulat"
)

type CollectionController struct {
	*chi.Mux
}

func NewCollectionController() *CollectionController {
	r := &CollectionController{
		Mux: chi.NewRouter(),
	}

	r.Get("/", wrapHandler(r.getCollections))
	r.With(validateRequest[sulat.Collection]()).Post("/", wrapHandler(r.createCollection))
	r.Delete("/", wrapHandler(r.removeSite))
	r.Route("/{collectionId}", func(sr chi.Router) {
		sr.Use(getCollectionCtx)
		sr.Get("/", wrapHandler(r.getCollection))
		sr.Delete("/", wrapHandler(r.removeCollection))
		sr.Route("/schema", func(sr chi.Router) {
			sr.Get("/", wrapHandler(r.getSchema))
			sr.Patch("/", wrapHandler(r.updateSchema))
		})
		sr.Mount("/records", NewRecordController())
	})

	return r
}

func (c *CollectionController) getCollections(w http.ResponseWriter, r *http.Request) error {
	site := getCurrentSite(r)
	collections, err := site.Collections()
	if err != nil {
		return err
	}
	return returnJson(w, collections)
}

func (c *CollectionController) createCollection(w http.ResponseWriter, r *http.Request) error {
	site := getCurrentSite(r)
	inst := getCurrentInstance(r)
	collection, err := getValidatedPayload[sulat.Collection](r)
	if err != nil {
		return err
	}

	collection.AttachSite(site)
	dataSource, err := inst.FindDataSource(collection.SourceId)
	if err != nil {
		return err
	}

	collection.Source = dataSource
	codec, err := inst.FindCodec(collection.CodecId)
	if err != nil {
		return err
	}

	collection.Codec = codec
	if _, err := site.CreateCollection(*collection); err != nil {
		return err
	}

	return returnJson(w, collection)
}

func (c *CollectionController) removeSite(w http.ResponseWriter, r *http.Request) error {
	site := getCurrentSite(r)
	inst := getCurrentInstance(r)
	if err := inst.RemoveSite(site.Id); err != nil {
		return err
	}
	return returnJson(w, nil)
}

func (c *CollectionController) getCollection(w http.ResponseWriter, r *http.Request) error {
	collection := getCurrentCollection(r)
	return returnJson(w, collection)
}

func (c *CollectionController) removeCollection(w http.ResponseWriter, r *http.Request) error {
	collection := getCurrentCollection(r)
	if err := collection.Site().RemoveCollection(collection.Id); err != nil {
		return err
	}
	return returnJson(w, nil)
}

func (c *CollectionController) getSchema(w http.ResponseWriter, r *http.Request) error {
	collection := getCurrentCollection(r)
	return returnJson(w, collection.Schema)
}

func (c *CollectionController) updateSchema(w http.ResponseWriter, r *http.Request) error {
	collection := getCurrentCollection(r)
	schema := []sulat.SchemaField{}
	if err := json.NewDecoder(r.Body).Decode(&schema); err != nil {
		return err
	}
	collection.Schema = schema
	return returnJson(w, collection.Schema)
}
