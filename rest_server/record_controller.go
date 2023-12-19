package rest_server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nedpals/sulatcms/sulat/query"
)

type RecordController struct {
	*chi.Mux
}

func NewRecordController() *RecordController {
	r := &RecordController{
		Mux: chi.NewRouter(),
	}

	r.With(getQueryCtx).Get("/", wrapHandler(r.getRecords))
	r.With(validateRecord).Post("/", wrapHandler(r.createRecord))
	r.With(getRecordCtx).Delete("/{recordId}", wrapHandler(r.deleteRecord))
	r.With(validateRecord).Patch("/{recordId}", wrapHandler(r.updateRecord))

	return r
}

func (rc *RecordController) getRecords(w http.ResponseWriter, r *http.Request) error {
	collection := getCurrentCollection(r)
	query := getDecodedQuery(r)

	records, err := collection.Find(query, nil)
	if err != nil {
		return err
	}
	return returnJson(w, records)
}

func (rc *RecordController) createRecord(w http.ResponseWriter, r *http.Request) error {
	collection := getCurrentCollection(r)
	record := getCurrentRecord(r)
	if err := collection.Insert(record, nil); err != nil {
		return err
	}
	return returnJson(w, record)
}

func (rc *RecordController) deleteRecord(w http.ResponseWriter, r *http.Request) error {
	collection := getCurrentCollection(r)
	record := getCurrentRecord(r)

	if err := collection.Delete(query.Eq("id", record.Id), nil); err != nil {
		return err
	}
	return returnJson(w, nil)
}

func (rc *RecordController) updateRecord(w http.ResponseWriter, r *http.Request) error {
	collection := getCurrentCollection(r)
	record := getCurrentRecord(r)
	if err := collection.Update(record, nil); err != nil {
		return err
	}
	return returnJson(w, record)
}
