package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nedpals/sulatcms/sulat"
	"github.com/nedpals/sulatcms/sulat/query"
)

func returnJson(wr http.ResponseWriter, data any) error {
	wr.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(wr).Encode(map[string]any{
		"data": data,
	})
}

type HandlerFunc func(http.ResponseWriter, *http.Request) error

func wrapHandler(fn func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			statusCode := http.StatusInternalServerError
			errorMessage := err.Error()

			if handlerErr, ok := err.(*sulat.ResponseError); ok {
				statusCode = handlerErr.StatusCode
				errorMessage = handlerErr.Message
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(statusCode)
			json.NewEncoder(w).Encode(map[string]any{
				"error": map[string]any{
					"status":  statusCode,
					"message": errorMessage,
				},
			})
		}
	}
}

type currentInstanceCtx struct{}

func getInstanceCtx(inst *sulat.Instance) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return wrapHandler(func(w http.ResponseWriter, r *http.Request) error {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), currentInstanceCtx{}, inst)))
			return nil
		})
	}
}

func getCurrentInstance(r *http.Request) *sulat.Instance {
	return r.Context().Value(currentInstanceCtx{}).(*sulat.Instance)
}

type currentSiteCtx struct{}

func getSiteCtx(next http.Handler) http.Handler {
	return wrapHandler(func(w http.ResponseWriter, r *http.Request) error {
		inst := getCurrentInstance(r)
		siteId := chi.URLParam(r, "siteId")

		if len(siteId) == 0 {
			siteIdHeader := r.Header.Get("X-Site-Id")
			if len(siteIdHeader) == 0 {
				return sulat.NewResponseError(http.StatusBadRequest, "site id is required")
			}

			siteId = siteIdHeader
		}

		site, err := inst.FindSite(siteId)
		if err != nil {
			return err
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), currentSiteCtx{}, site)))
		return nil
	})
}

func getCurrentSite(r *http.Request) *sulat.Site {
	return r.Context().Value(currentSiteCtx{}).(*sulat.Site)
}

type currentCollectionCtx struct{}

func getCollectionCtx(next http.Handler) http.Handler {
	return wrapHandler(func(w http.ResponseWriter, r *http.Request) error {
		site := getCurrentSite(r)
		collectionId := chi.URLParam(r, "collectionId")
		collection, err := site.FindCollection(collectionId)
		if err != nil {
			return err
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), currentCollectionCtx{}, collection)))
		return nil
	})
}

func getCurrentCollection(r *http.Request) *sulat.Collection {
	return r.Context().Value(currentCollectionCtx{}).(*sulat.Collection)
}

type currentDataSourceCtx struct{}

func getDataSourceCtx(next http.Handler) http.Handler {
	return wrapHandler(func(w http.ResponseWriter, r *http.Request) error {
		inst := getCurrentInstance(r)
		dataSourceId := chi.URLParam(r, "dataSourceId")
		dataSource, err := inst.FindDataSource(dataSourceId)
		if err != nil {
			return err
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), currentDataSourceCtx{}, dataSource)))
		return nil
	})
}

type currentRecordCtx struct{}

func getRecordCtx(next http.Handler) http.Handler {
	return wrapHandler(func(w http.ResponseWriter, r *http.Request) error {
		collection := getCurrentCollection(r)
		recordId := chi.URLParam(r, "recordId")
		record, err := collection.Source.Get(collection.Id, recordId, nil)
		if err != nil {
			return err
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), currentRecordCtx{}, record)))
		return nil
	})
}

type currentQueryCtx struct{}

func getQueryCtx(next http.Handler) http.Handler {
	return wrapHandler(func(w http.ResponseWriter, r *http.Request) error {
		query, err := query.ParseFromRequest(r)
		if err != nil {
			return err
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), currentQueryCtx{}, query)))
		return nil
	})
}

func getDecodedQuery(r *http.Request) *query.Query {
	return r.Context().Value(currentQueryCtx{}).(*query.Query)
}

func getCurrentRecord(r *http.Request) *sulat.Record {
	return r.Context().Value(currentRecordCtx{}).(*sulat.Record)
}

func getCurrentDataSource(r *http.Request) *sulat.DataSource {
	return r.Context().Value(currentDataSourceCtx{}).(*sulat.DataSource)
}

func validateRecord(next http.Handler) http.Handler {
	return wrapHandler(func(w http.ResponseWriter, r *http.Request) error {
		collection := getCurrentCollection(r)
		payload := map[string]any{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			return err
		}

		if err := collection.Schema.Validate(payload); err != nil {
			return err
		}

		record := &sulat.Record{
			Id:         payload["id"].(string),
			Data:       payload,
			Collection: collection,
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), currentRecordCtx{}, record)))
		return nil
	})
}
