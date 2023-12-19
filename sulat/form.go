package sulat

import (
	"database/sql/driver"
)

type FormSchema []FormBlock

func (f *FormSchema) Scan(src any) error {
	return scanJson(src, f, "FormSchema")
}

func (f FormSchema) Value() (driver.Value, error) {
	return driverValueJson(f)
}

type FormLocation struct {
	Name       string         `json:"name"`
	Label      string         `json:"label"`
	Blocks     []FormBlock    `json:"blocks"`
	Properties map[string]any `json:"properties"`
}

type FormBlock struct {
	Field      string         `json:"field"`
	Type       string         `json:"type"`
	Location   string         `json:"location"`
	Properties map[string]any `json:"properties"`
}
