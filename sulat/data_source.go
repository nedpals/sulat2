package sulat

import (
	"errors"
	"net/http"
	"path/filepath"
	"slices"

	"github.com/jmoiron/sqlx"
	"github.com/nedpals/sulatcms/sulat/query"
	"github.com/spf13/afero"
	"golang.org/x/exp/maps"
)

type DataSource struct {
	instance           *Instance
	Id                 string         `json:"id" db:"id" mapstructure:"id"`
	Name               string         `json:"name" db:"name" mapstructure:"name"`
	Config             map[string]any `json:"config" db:"config" mapstructure:"config,omitempty"`
	ProviderId         string         `json:"provider" db:"provider" mapstructure:"provider"`
	DataSourceProvider `json:"-" db:"-"`
}

func (ds DataSource) ValidationSchema() Schema {
	return Schema{
		StringSchemaField{
			BaseField: BaseField{
				FieldName: "id",
				Required:  true,
			},
		},
		StringSchemaField{
			BaseField: BaseField{
				FieldName: "name",
				Required:  true,
			},
		},
		StringSchemaField{
			BaseField: BaseField{
				FieldName: "provider",
				Required:  true,
			},
		},
	}
}

func (ds *DataSource) Initialize() error {
	if ds.DataSourceProvider == nil {
		provider, err := ds.instance.FindDataSourceProvider(ds.ProviderId)
		if err != nil {
			return err
		}

		ds.DataSourceProvider = provider
	}

	var err error
	ds.DataSourceProvider, err = ds.DataSourceProvider.WithConfig(ds.Config)
	if err != nil {
		return err
	}

	return ds.DataSourceProvider.Initialize(ds.instance)
}

func fetchDataSources(dataSources *[]*DataSource, db *sqlx.DB) error {
	return db.Select(dataSources, "SELECT * FROM data_sources")
}

func createDataSource(dataSource *DataSource, db *sqlx.DB) error {
	_, err := db.NamedExec("INSERT INTO data_sources (id, name, config, provider) VALUES (:id, :name, :config, :provider)", dataSource)
	return err
}

func updateDataSource(dataSource *DataSource, db *sqlx.DB) error {
	_, err := db.NamedExec("UPDATE data_sources SET name = :name, config = :config, provider = :provider WHERE id = :id", dataSource)
	return err
}

func removeDataSource(dataSourceId string, db *sqlx.DB) error {
	_, err := db.NamedExec("DELETE FROM data_sources WHERE id = :id", dataSourceId)
	return err
}

type DataSourceProvider interface {
	Initialize(*Instance) error
	Properties() DataSourceProviderProperties
	WithConfig(config map[string]any) (DataSourceProvider, error)
	Get(collectionId string, id string, opts map[string]any) (*Record, error)
	Find(collectionId string, query *query.Query, opts map[string]any) ([]*Record, error)
	Insert(collectionId string, record *Record, opts map[string]any) error
	Update(collectionId string, record *Record, opts map[string]any) error
	Delete(collectionId string, query *query.Query, opts map[string]any) error
}

type DataSourceProviderProperties struct {
	Id           string
	Name         string
	Version      string
	Config       map[string]any
	ConfigSchema Schema
}

// DATA SOURCE PROVIDER IMPLEMENTATIONS

// FileDataSourceProvider is a data source provider that uses the file system as its data source
type FileDataSourceProvider struct {
	// FS is the file system to use for this provider
	FS afero.Fs

	// Codecs is the codecs registry to use for this provider (should only be used in testing as possible)
	codecs CodecRegistry

	// ConfigPath is the path where sulat.toml is located. This will also be used to determine the root directory
	ConfigPath string

	// Root is the root directory of the data source
	Root string

	// Collections is a map of collection ids to globbed paths
	Collections map[string]string

	// cachedCollections is a map of collection ids to collections
	cachedCollections map[string]*Collection

	// records is a map of collection ids to records
	records map[string]map[string]*Record
}

func injectConfigToProvider(p *FileDataSourceProvider, config map[string]any) {
	if len(p.ConfigPath) == 0 {
		if newConfigPath, ok := config["config_path"]; ok {
			if newConfigPathStr, ok := newConfigPath.(string); ok {
				p.ConfigPath = newConfigPathStr
			}
		}
	}

	if rawCollections, ok := config["collections"]; ok {
		if p.Collections == nil {
			p.Collections = make(map[string]string)
		}

		if collections, ok := rawCollections.(map[string]string); ok {
			maps.Copy(p.Collections, collections)
		} else if collections, ok := rawCollections.(map[string]any); ok {
			for collectionId, glob := range collections {
				glob, ok := glob.(string)
				if !ok {
					continue
				}
				p.Collections[collectionId] = glob
			}
		}
	}

	if rawRoot, ok := config["root"]; ok {
		if root, ok := rawRoot.(string); ok {
			p.Root = root
		}
	}
}

func (p *FileDataSourceProvider) Initialize(i *Instance) error {
	if !slices.Equal(p.codecs, i.Codecs()) {
		p.codecs = i.Codecs()
	}

	if p.FS == nil {
		p.FS = afero.NewOsFs()
	}

	if p.cachedCollections == nil {
		p.cachedCollections = make(map[string]*Collection)
	}

	var importErrors []error

	for collectionId, glob := range p.Collections {
		// create collection first
		collection := &Collection{
			Id: collectionId,
		}

		// TODO: replace this and make it "importable" to site instead
		p.cachedCollections[collectionId] = collection

		// import records
		files, err := afero.Glob(p.FS, filepath.Join(p.Root, glob))
		if err != nil {
			return err
		}

		if len(files) == 0 {
			continue
		}

		records := map[string]*Record{}
		for _, filename := range files {
			file, err := p.FS.Open(filename)
			if err != nil {
				importErrors = append(importErrors, err)
				continue
			}

			stat, err := file.Stat()
			if err != nil {
				importErrors = append(importErrors, err)
				continue
			}

			codec, err := i.Codecs().FindByFileName(stat.Name())
			if err != nil {
				importErrors = append(importErrors, err)
				continue
			}

			record, err := codec.Deserialize(stat.Name(), file)
			if err != nil {
				importErrors = append(importErrors, err)
				continue
			}

			// strip root path from filename
			finalFilename := filename
			relFilename, err := filepath.Rel(p.Root, filename)
			if err == nil {
				finalFilename = relFilename
			}

			records[finalFilename] = &Record{
				Id:         record.Id,
				Collection: collection,
				Data:       record.Data,
			}
		}

		if p.records == nil {
			p.records = make(map[string]map[string]*Record)
		}

		p.records[collectionId] = records
	}

	return errors.Join(importErrors...)
}

func (p *FileDataSourceProvider) Properties() DataSourceProviderProperties {
	return DataSourceProviderProperties{
		Id:      "fs",
		Name:    "File System",
		Version: "1.0.0",
		ConfigSchema: Schema{
			StringSchemaField{
				BaseField: BaseField{
					FieldName:  "config_path",
					FieldLabel: "Config path",
				},
			},
			StringSchemaField{
				BaseField: BaseField{
					FieldName:  "root",
					FieldLabel: "Root",
				},
			},
			KVGroupSchemaField{
				BaseField: BaseField{
					FieldName:  "collections",
					FieldLabel: "Collections",
				},
				KeySchema: StringSchemaField{
					BaseField: BaseField{
						FieldName:  "collection_id",
						FieldLabel: "Collection ID",
						Required:   true,
					},
				},
				ValueSchema: StringSchemaField{
					BaseField: BaseField{
						FieldName:  "path",
						FieldLabel: "Collection path",
						Required:   true,
					},
				},
			},
		},
	}
}

func (p *FileDataSourceProvider) WithConfig(config map[string]any) (DataSourceProvider, error) {
	newProvider := &FileDataSourceProvider{
		FS: p.FS,
	}

	injectConfigToProvider(newProvider, config)

	if len(newProvider.Root) != 0 && len(newProvider.Collections) != 0 {
		return newProvider, nil
	}

	// import config file if no root or collections are specified

	// get nearest config file if no config path is specified
	if len(newProvider.ConfigPath) == 0 && len(newProvider.Root) != 0 {
		configFilePaths, err := afero.Glob(newProvider.FS, filepath.Join(newProvider.Root, "**", "sulat.toml"))
		if err != nil {
			return nil, err
		} else if len(configFilePaths) == 0 {
			return nil, errors.New("no sulat.toml found")
		}

		newProvider.ConfigPath = configFilePaths[0]
	}

	// still if there's no config path, return an error
	if len(newProvider.ConfigPath) == 0 {
		return nil, errors.New("no config path found. create a sulat.toml file or specify a config path")
	}

	configFromFile, err := ParseConfigFile(newProvider.FS, newProvider.ConfigPath)
	if err != nil {
		return nil, err
	}

	// validate config with config schema first
	schema := newProvider.Properties().ConfigSchema
	if err := schema.Validate(configFromFile); err != nil {
		return nil, err
	}

	injectConfigToProvider(newProvider, configFromFile)
	return newProvider, nil
}

func (p *FileDataSourceProvider) fetchRecord(collectionId string, id string) (string, *Record, error) {
	records, found := p.records[collectionId]
	if !found {
		return "", nil, NewResponseError(http.StatusNotFound, "collection not found")
	}

	for filename, record := range records {
		if record.Id == id {
			return filename, record, nil
		}
	}

	return "", nil, NewResponseError(http.StatusNotFound, "record not found")
}

func (p *FileDataSourceProvider) saveRecord(filename string, record *Record) error {
	if len(filename) == 0 {
		return NewResponseError(http.StatusBadRequest, "fileName is required")
	}

	data, err := record.Serialize()
	if err != nil {
		return err
	}

	if err := afero.WriteFile(p.FS, filepath.Join(p.Root, filename), data, 0644); err != nil {
		return err
	}

	p.records[record.Collection.Id][filename] = record
	return nil
}

func (p *FileDataSourceProvider) Get(collectionId string, id string, opts map[string]any) (*Record, error) {
	_, record, err := p.fetchRecord(collectionId, id)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (p *FileDataSourceProvider) Find(collectionId string, query *query.Query, opts map[string]any) ([]*Record, error) {
	records, collectionFound := p.records[collectionId]
	if !collectionFound {
		return nil, NewResponseError(http.StatusNotFound, "collection not found")
	}

	if query == nil {
		return maps.Values(records), nil
	}

	// match query against the records
	found := []*Record{}

	for _, record := range records {
		if query.Match(record) {
			found = append(found, record)
		}
	}

	if len(found) == 0 {
		return nil, NewResponseError(http.StatusNotFound, "no records found")
	}

	return found, nil
}

func (p *FileDataSourceProvider) Insert(collectionId string, record *Record, opts map[string]any) error {
	// check for duplicate id
	_, existingRecord, _ := p.fetchRecord(collectionId, record.Id)
	if existingRecord != nil {
		return NewResponseError(http.StatusConflict, "record already exists")
	}

	if record.Collection == nil {
		collection, found := p.cachedCollections[collectionId]
		if !found {
			return NewResponseError(http.StatusNotFound, "collection not found")
		}

		record.Collection = collection
	}

	if record.Codec == nil {
		codec, err := p.codecs.FindByFileName(record.Id)
		if err != nil {
			return err
		}

		record.Codec = codec
	}

	// save to fs
	filename := filepath.Join(collectionId, record.Id)
	return p.saveRecord(filename, record)
}

func (p *FileDataSourceProvider) Update(collectionId string, updateRecord *Record, opts map[string]any) error {
	filename, record, err := p.fetchRecord(collectionId, updateRecord.Id)
	if err != nil {
		return err
	}

	// merge record with updateRecord
	updatedRecord := &Record{
		Id:         record.Id,
		Collection: record.Collection,
		Data:       updateRecord.Data,
		Codec:      record.Codec,
	}

	return p.saveRecord(filename, updatedRecord)
}

func (p *FileDataSourceProvider) Delete(collectionId string, query *query.Query, opts map[string]any) error {
	records, collectionFound := p.records[collectionId]
	if !collectionFound {
		return NewResponseError(http.StatusNotFound, "collection not found")
	} else if query == nil {
		return NewResponseError(http.StatusBadRequest, "query is required")
	}

	// match query against the records
	maps.DeleteFunc(records, func(k string, r *Record) bool {
		isMatched := query.Match(r)
		if isMatched {
			p.FS.Remove(k)
		}
		return isMatched
	})

	return nil
}
