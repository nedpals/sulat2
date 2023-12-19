package sulat

import (
	"database/sql"
	"net/http"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/afero"
	"golang.org/x/exp/slices"
)

type Site struct {
	instance            *Instance
	Id                  string         `json:"id" db:"id"`
	Name                string         `json:"name" db:"name"`
	Other               map[string]any `json:"other" db:"other"`
	collections         []*Collection
	DefaultDataSourceId sql.NullString `json:"default_data_source" db:"default_data_source"`
	DefaultDataSource   DataSource     `json:"-" db:"-"`
}

func fetchSites(sites *[]*Site, db *sqlx.DB) error {
	return db.Select(sites, "SELECT * FROM sites")
}

func createSite(site *Site, db *sqlx.DB) error {
	_, err := db.NamedExec("INSERT INTO sites (id, name, default_data_source) VALUES (:id, :name, :default_data_source)", site)
	return err
}

func updateSite(site *Site, db *sqlx.DB) error {
	_, err := db.NamedExec("UPDATE sites SET name = :name, default_data_source = :default_data_source WHERE id = :id", site)
	return err
}

func removeSite(siteId string, db *sqlx.DB) error {
	_, err := db.NamedExec("DELETE FROM sites WHERE id = :id", siteId)
	return err
}

func (s *Site) NewCollection() *Collection {
	collection := &Collection{}
	return s.attachCollection(collection)
}

func (s *Site) attachCollection(c *Collection) *Collection {
	if c.site == nil {
		c.AttachSite(s)
	}
	return c
}

func (s *Site) Collections() ([]*Collection, error) {
	if s.collections == nil {
		if err := fetchCollections(&s.collections, s.instance.db); err != nil {
			return nil, err
		}

		for _, c := range s.collections {
			s.attachCollection(c)
		}
	}
	return s.collections, nil
}

func (s *Site) FindCollection(collectionId string) (*Collection, error) {
	collections, err := s.Collections()
	if err != nil {
		return nil, err
	}

	for _, collection := range collections {
		if collection.Id == collectionId {
			return s.attachCollection(collection), nil
		}
	}

	return nil, NewResponseError(http.StatusNotFound, "collection not found")
}

func (s *Site) CreateCollection(c Collection) (*Collection, error) {
	collection := &Collection{
		site:     s,
		Id:       c.Id,
		Name:     c.Name,
		SourceId: c.Source.Properties().Id,
		Source:   c.Source,
	}

	if err := createCollection(collection, s.instance.db); err != nil {
		return nil, err
	}

	s.collections = append(s.collections, collection)
	return collection, nil
}

func (s *Site) RemoveCollection(collectionId string) error {
	collection, err := s.FindCollection(collectionId)
	if err != nil {
		return err
	}

	if err := removeCollection(collection.Id, s.instance.db); err != nil {
		return err
	}

	s.collections = slices.DeleteFunc(s.collections, func(c *Collection) bool {
		return c.Id == collectionId
	})
	return nil
}

func (s *Site) UpdateCollection(collection *Collection) error {
	if err := updateCollection(collection, s.instance.db); err != nil {
		return err
	}

	for i, c := range s.collections {
		if c.Id == collection.Id {
			s.collections[i] = collection
			break
		}
	}

	return nil
}

func ParseConfigFile(fs afero.Fs, path string) (map[string]any, error) {
	configFile, err := afero.ReadFile(fs, path)
	if err != nil {
		return nil, err
	}

	var config map[string]any
	if err := toml.Unmarshal(configFile, &config); err != nil {
		return nil, err
	}

	config["root"] = filepath.Dir(path)
	return config, nil
}
