package sulat

import (
	"github.com/jmoiron/sqlx"
	"github.com/nedpals/sulatcms/sulat/query"
)

// Collection represents a collection of records
type Collection struct {
	site       *Site
	Id         string         `json:"id" db:"id"`
	Name       string         `json:"name" db:"name"`
	Metadata   map[string]any `json:"-" db:"metadata"`
	Source     *DataSource    `json:"-" db:"-"`
	SourceId   string         `json:"source" db:"source"`
	Schema     Schema         `json:"-" db:"schema"`
	FormSchema FormSchema     `json:"form_schema" db:"form_schema"`
}

// Site gets the site associated with the collection
func (c *Collection) Site() *Site {
	return c.site
}

// Get gets a record from the collection
func (c *Collection) Get(id string, opts map[string]any) (*Record, error) {
	return c.Source.Get(c.Id, id, opts)
}

// Find finds records from the collection
func (c *Collection) Find(query *query.Query, opts map[string]any) ([]*Record, error) {
	return c.Source.Find(c.Id, query, opts)
}

// Insert inserts a record into the collection
func (c *Collection) Insert(record *Record, opts map[string]any) error {
	return c.Source.Insert(c.Id, record, opts)
}

// Update updates a record from the collection
func (c *Collection) Update(record *Record, opts map[string]any) error {
	return c.Source.Update(c.Id, record, opts)
}

// Delete deletes a record from the collection
func (c *Collection) Delete(query *query.Query, opts map[string]any) error {
	return c.Source.Delete(c.Id, query, opts)
}

func fetchCollections(collections *[]*Collection, db *sqlx.DB) error {
	return db.Select(collections, "SELECT * FROM collections")
}

func createCollection(collection *Collection, db *sqlx.DB) error {
	_, err := db.NamedExec("INSERT INTO collections (id, label, source) VALUES (:id, :label, :source)", collection)
	return err
}

func updateCollection(collection *Collection, db *sqlx.DB) error {
	_, err := db.NamedExec("UPDATE collections SET label = :label, source = :source WHERE id = :id", collection)
	return err
}

func removeCollection(collectionId string, db *sqlx.DB) error {
	_, err := db.NamedExec("DELETE FROM collections WHERE id = :id", collectionId)
	return err
}
