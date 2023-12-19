package sulat

import (
	"strconv"
	"strings"
)

type fieldParser struct {
	idx  int
	keys []string
	data any
}

func (fp *fieldParser) get() any {
	if fp.idx == len(fp.keys) {
		return fp.data
	} else if fp.idx > len(fp.keys) || fp.data == nil {
		return nil
	}

	key := fp.keys[fp.idx]

	switch v := fp.data.(type) {
	case map[string]any:
		vv, ok := v[key]
		if !ok {
			return nil
		}
		fp.data = vv
		fp.idx++
		return fp.get()
	case []any:
		idx, err := strconv.Atoi(key)
		if err != nil || idx < 0 || idx >= len(v) {
			return nil
		}

		fp.data = v[idx]
		fp.idx++
		return fp.get()
	default:
		return fp.data
	}
}

func parseField(field string, data map[string]any) *fieldParser {
	return &fieldParser{
		keys: strings.Split(field, "."),
		data: data,
	}
}

// TODO: add tags support
type Record struct {
	Id         string
	Data       map[string]any
	Collection *Collection
}

// Get returns the value of a field
func (r *Record) Get(field string) any {
	if field == "id" {
		return r.Id
	}

	fp := parseField(field, r.Data)
	return fp.get()
}

// Set sets the value of a field
func (r *Record) Title() string {
	return r.Get("title").(string)
}
