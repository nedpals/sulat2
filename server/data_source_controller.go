package server

import (
	"encoding/json"
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
	r.Post("/", wrapHandler(r.createDataSource))
	r.Route("/{dataSourceId}", func(sr chi.Router) {
		sr.Use(getDataSourceCtx)
		sr.Get("/", wrapHandler(r.getDataSource))
		sr.Delete("/", wrapHandler(r.removeDataSource))
		sr.Patch("/", wrapHandler(r.updateDataSource))
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
	dataSource := inst.NewDataSource()
	if err := json.NewDecoder(r.Body).Decode(&dataSource); err != nil {
		return err
	}

	if len(dataSource.Id) == 0 {
		return sulat.NewResponseError(http.StatusBadRequest, "data source id is required")
	}

	if len(dataSource.Name) == 0 {
		return sulat.NewResponseError(http.StatusBadRequest, "data source name is required")
	}

	if len(dataSource.Config) == 0 {
		return sulat.NewResponseError(http.StatusBadRequest, "data source config is required")
	}

	if len(dataSource.ProviderId) == 0 {
		return sulat.NewResponseError(http.StatusBadRequest, "data source provider is required")
	}

	if _, err := inst.CreateDataSource(dataSource.Id, dataSource.Name, dataSource.Config, dataSource.ProviderId); err != nil {
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
	if err := json.NewDecoder(r.Body).Decode(&dataSource); err != nil {
		return err
	}

	if len(dataSource.Name) == 0 {
		return sulat.NewResponseError(http.StatusBadRequest, "data source name is required")
	}

	if len(dataSource.Config) == 0 {
		return sulat.NewResponseError(http.StatusBadRequest, "data source config is required")
	}

	if err := inst.UpdateDataSource(dataSource); err != nil {
		return err
	}

	return returnJson(w, dataSource)
}
