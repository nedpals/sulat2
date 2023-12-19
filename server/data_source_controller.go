package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nedpals/sulatcms/sulat"
)

type DataSourceController struct {
	*chi.Mux
}

func NewDataSourceController() *DataSourceController {
	r := &DataSourceController{
		Mux: chi.NewRouter(),
	}

	r.Get("/providers", wrapHandler(r.getProviders))
	r.Get("/", wrapHandler(r.getDataSources))
	r.With(validateRequest[sulat.DataSource]()).Post("/", wrapHandler(r.createDataSource))
	r.Route("/{dataSourceId}", func(sr chi.Router) {
		sr.Use(getDataSourceCtx)
		sr.Get("/", wrapHandler(r.getDataSource))
		sr.Delete("/", wrapHandler(r.removeDataSource))
		sr.With(validateRequest[sulat.DataSource]()).Patch("/", wrapHandler(r.updateDataSource))
	})

	return r
}

func (c *DataSourceController) getProviders(w http.ResponseWriter, r *http.Request) error {
	inst := getCurrentInstance(r)
	dataSourceProviders := inst.DataSourceProviders()
	return returnJson(w, dataSourceProviders)
}

func (c *DataSourceController) getDataSources(w http.ResponseWriter, r *http.Request) error {
	inst := getCurrentInstance(r)
	dataSources, err := inst.DataSources()
	if err != nil {
		return err
	}
	return returnJson(w, dataSources)
}

func (c *DataSourceController) createDataSource(w http.ResponseWriter, r *http.Request) error {
	inst := getCurrentInstance(r)
	dataSource, err := getValidatedPayload[sulat.DataSource](r)
	if err != nil {
		return err
	}

	if _, err := inst.CreateDataSource(*dataSource); err != nil {
		return err
	}

	return returnJson(w, dataSource)
}

func (c *DataSourceController) getDataSource(w http.ResponseWriter, r *http.Request) error {
	dataSource := getCurrentDataSource(r)
	return returnJson(w, dataSource)
}

func (c *DataSourceController) removeDataSource(w http.ResponseWriter, r *http.Request) error {
	inst := getCurrentInstance(r)
	dataSource := getCurrentDataSource(r)
	if err := inst.RemoveDataSource(dataSource.Id); err != nil {
		return err
	}
	return returnJson(w, nil)
}

func (c *DataSourceController) updateDataSource(w http.ResponseWriter, r *http.Request) error {
	dataSource := getCurrentDataSource(r)
	inst := getCurrentInstance(r)
	validated, err := getValidatedPayload[sulat.DataSource](r)
	if err != nil {
		return err
	}

	dataSource.Id = validated.Id
	dataSource.Name = validated.Name
	if err := inst.UpdateDataSource(dataSource); err != nil {
		return err
	}

	return returnJson(w, dataSource)
}
