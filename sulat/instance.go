package sulat

import (
	_ "embed"
	"net/http"

	"github.com/jmoiron/sqlx"
	"golang.org/x/exp/slices"
	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var dbSchema string

type Instance struct {
	db                  *sqlx.DB
	sites               []*Site
	dataSourceProviders []DataSourceProvider
	dataSources         []*DataSource
	codecs              CodecRegistry
}

// NewInstance creates a new instance
func NewInstance(dbLocation string) (*Instance, error) {
	if len(dbLocation) == 0 {
		dbLocation = ":memory:"
	}

	db, err := sqlx.Open("sqlite", dbLocation)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(dbSchema); err != nil {
		return nil, err
	}

	inst := &Instance{
		db: db,
	}

	if err := fetchDataSources(&inst.dataSources, inst.db); err != nil {
		return nil, err
	}

	if err := fetchSites(&inst.sites, inst.db); err != nil {
		return nil, err
	}

	return inst, nil
}

// Sites returns all sites
func (i *Instance) Sites() ([]*Site, error) {
	if i.sites == nil {
		if err := fetchSites(&i.sites, i.db); err != nil {
			return nil, err
		}
	}
	return i.sites, nil
}

// FindSite finds a site by id
func (i *Instance) FindSite(id string) (*Site, error) {
	sites, err := i.Sites()
	if err != nil {
		return nil, err
	}

	for _, site := range sites {
		if site.Id == id {
			return site, nil
		}
	}

	return nil, NewResponseError(http.StatusNotFound, "site not found")
}

// RemoveSite removes a site
func (i *Instance) RemoveSite(siteId string) error {
	site, err := i.FindSite(siteId)
	if err != nil {
		return err
	}

	if err := removeSite(site.Id, i.db); err != nil {
		return err
	}

	i.sites = slices.DeleteFunc(i.sites, func(s *Site) bool {
		return s.Id == siteId
	})
	return nil
}

type CreateSiteParams struct {
	Name       string
	DataSource DataSource
}

// CreateSite creates a new site
func (i *Instance) CreateSite(id string, params CreateSiteParams) (*Site, error) {
	siteName := id
	if len(params.Name) != 0 {
		siteName = params.Name
	}

	site := &Site{
		instance:          i,
		Id:                id,
		Name:              siteName,
		DefaultDataSource: params.DataSource,
	}

	if err := createSite(site, i.db); err != nil {
		return nil, err
	}

	i.sites = append(i.sites, site)
	return site, nil
}

// UpdateSite updates a site
func (i *Instance) UpdateSite(site *Site) error {
	if err := updateSite(site, i.db); err != nil {
		return err
	}

	for idx, s := range i.sites {
		if s.Id == site.Id {
			i.sites[idx] = site
			break
		}
	}

	return nil
}

// RegisterDataSourceProvider registers a data source provider
func (i *Instance) RegisterDataSourceProvider(dataSource DataSourceProvider) {
	if i.dataSourceProviders == nil {
		i.dataSourceProviders = []DataSourceProvider{}
	}
	i.dataSourceProviders = append(i.dataSourceProviders, dataSource)
}

// FindDataSourceProvider finds a data source provider by id
func (i *Instance) FindDataSourceProvider(id string) (DataSourceProvider, error) {
	for _, dataSource := range i.dataSourceProviders {
		if dataSource.Properties().Id == id {
			return dataSource, nil
		}
	}
	return nil, NewResponseError(http.StatusNotFound, "data source not found")
}

// DataSourceProviders returns all data source providers
func (i *Instance) DataSourceProviders() []DataSourceProvider {
	return i.dataSourceProviders
}

// DataSources returns all data sources
func (i *Instance) DataSources() ([]*DataSource, error) {
	if i.dataSources == nil {
		if err := fetchDataSources(&i.dataSources, i.db); err != nil {
			return nil, err
		}

		for _, dataSource := range i.dataSources {
			provider, err := i.FindDataSourceProvider(dataSource.ProviderId)
			if err != nil {
				// TODO: accumulate errors
				continue
			}
			dataSource.DataSourceProvider = provider
		}
	}
	return i.dataSources, nil
}

// attachDataSource attaches a data source to an instance
func (i *Instance) attachDataSource(dataSource *DataSource) *DataSource {
	if dataSource.instance == nil {
		dataSource.instance = i
	}
	return dataSource
}

// NewDataSource creates a new data source
func (i *Instance) NewDataSource(id, name string, provider DataSourceProvider, config map[string]any) *DataSource {
	ds := &DataSource{
		Id:                 id,
		Name:               name,
		Config:             config,
		ProviderId:         provider.Properties().Id,
		DataSourceProvider: provider,
	}

	ds = i.attachDataSource(ds)
	schema := provider.Properties().ConfigSchema
	if err := schema.Validate(config); err != nil {
		panic(err)
	}

	if err := ds.Initialize(); err != nil {
		panic(err)
	}

	return ds
}

// FindDataSource finds a data source by id
func (i *Instance) FindDataSource(id string) (*DataSource, error) {
	dataSources, err := i.DataSources()
	if err != nil {
		return nil, err
	}

	for _, dataSource := range dataSources {
		if dataSource.Id == id {
			return dataSource, nil
		}
	}

	return nil, NewResponseError(http.StatusNotFound, "data source not found")
}

// CreateDataSource creates a new data source
func (i *Instance) CreateDataSource(ds DataSource) (*DataSource, error) {
	provider, err := i.FindDataSourceProvider(ds.ProviderId)
	if err != nil {
		return nil, err
	}

	dataSource := &DataSource{
		instance:           i,
		Id:                 ds.Id,
		Name:               ds.Name,
		Config:             ds.Config,
		ProviderId:         ds.ProviderId,
		DataSourceProvider: provider,
	}

	if err := createDataSource(dataSource, i.db); err != nil {
		return nil, err
	}

	i.dataSources = append(i.dataSources, dataSource)
	return dataSource, nil
}

// RemoveDataSource removes a data source
func (i *Instance) RemoveDataSource(id string) error {
	// check for dependencies first before deleting
	sites, err := i.Sites()
	if err == nil {
		for _, site := range sites {
			if site.DefaultDataSourceId.Valid && site.DefaultDataSourceId.String == id {
				return NewResponseError(http.StatusBadRequest, "data source is in use")
			}

			collections, err := site.Collections()
			if err != nil {
				continue
			}

			// check if data source is in use by a collection
			for _, collection := range collections {
				if collection.SourceId == id {
					return NewResponseError(http.StatusBadRequest, "data source is in use")
				}
			}
		}
	}

	dataSource, err := i.FindDataSource(id)
	if err != nil {
		return err
	}

	if err := removeDataSource(dataSource.Id, i.db); err != nil {
		return err
	}

	i.dataSources = slices.DeleteFunc(i.dataSources, func(ds *DataSource) bool {
		return ds.Id == id
	})
	return nil
}

// UpdateDataSource updates a data source
func (i *Instance) UpdateDataSource(dataSource *DataSource) error {
	if err := updateDataSource(dataSource, i.db); err != nil {
		return err
	}

	for idx, ds := range i.dataSources {
		if ds.Id == dataSource.Id {
			i.dataSources[idx] = dataSource
			break
		}
	}

	return nil
}

// Codecs returns all codecs
func (i *Instance) Codecs() CodecRegistry {
	return i.codecs
}

// FindCodec finds a codec by id
func (i *Instance) FindCodec(id string) (*Codec, error) {
	return i.Codecs().Find(id)
}

// RegisterCodec registers a codec
func (i *Instance) RegisterCodec(codec *Codec) error {
	return i.codecs.Register(codec)
}

// RegisterCodecs registers multiple codecs
func (i *Instance) RegisterCodecs(codecs ...*Codec) error {
	return i.codecs.RegisterMultiple(codecs...)
}

// UpdateCodec updates a codec
func (i *Instance) UpdateCodec(codec *Codec) error {
	return i.codecs.Update(codec)
}

// RemoveCodec removes a codec
func (i *Instance) RemoveCodec(id string) error {
	return i.codecs.Remove(id)
}
