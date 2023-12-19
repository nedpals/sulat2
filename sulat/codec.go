package sulat

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"golang.org/x/exp/slices"
)

type Codec struct {
	Id             string
	FileExtensions []string
	ContentTypes   []string
	OnDeserialize  func(io.Reader) (map[string]any, error)
	OnSerialize    func(*Record) ([]byte, error)
}

func (c *Codec) Deserialize(id string, r io.Reader) (*Record, error) {
	data, err := c.OnDeserialize(r)
	if err != nil {
		return nil, err
	}

	return &Record{
		Id:    id,
		Data:  data,
		Codec: c,
	}, nil
}

func (c *Codec) Serialize(r *Record) ([]byte, error) {
	return c.OnSerialize(r)
}

type CodecRegistry []*Codec

// Register registers a codec
func (c *CodecRegistry) Register(codec *Codec) error {
	if *c == nil {
		*c = []*Codec{}
	}

	// check if codec already exists
	for _, c := range *c {
		if c.Id == codec.Id {
			return fmt.Errorf("codec %s already exists. UpdateCodec to override existing codec", codec.Id)
		}
	}

	*c = append(*c, codec)
	return nil
}

// RegisterMultiple registers multiple codecs
func (c *CodecRegistry) RegisterMultiple(codecs ...*Codec) error {
	for _, codec := range codecs {
		if err := c.Register(codec); err != nil {
			return err
		}
	}
	return nil
}

// Update updates a codec
func (c *CodecRegistry) Update(codec *Codec) error {
	if *c == nil {
		*c = []*Codec{}
	}

	// check if codec already exists
	for idx, cc := range *c {
		if cc.Id == codec.Id {
			(*c)[idx] = codec
			return nil
		}
	}

	return c.Register(codec)
}

// Remove removes a codec
func (c *CodecRegistry) Remove(id string) error {
	_, err := c.Find(id)
	if err != nil {
		return err
	}

	*c = slices.DeleteFunc(*c, func(c *Codec) bool {
		return c.Id == id
	})
	return nil
}

// Find finds a codec by id
func (c CodecRegistry) Find(id string) (*Codec, error) {
	for _, codec := range c {
		if codec.Id == id {
			return codec, nil
		}
	}
	return nil, NewResponseError(http.StatusNotFound, "codec not found")
}

// FindByFileExtension finds a codec by file extension
func (c CodecRegistry) FindByFileExtension(ext string) (*Codec, error) {
	for _, codec := range c {
		for _, codecExt := range codec.FileExtensions {
			if codecExt == ext {
				return codec, nil
			}
		}
	}
	return nil, NewResponseError(http.StatusNotFound, "codec not found")
}

// FindByFileName finds a codec by file name
func (c CodecRegistry) FindByFileName(fileName string) (*Codec, error) {
	ext := filepath.Ext(fileName)
	return c.FindByFileExtension(ext)
}

// FindByContentType finds a codec by content type
func (c CodecRegistry) FindByContentType(contentType string) (*Codec, error) {
	for _, codec := range c {
		for _, codecContentType := range codec.ContentTypes {
			if codecContentType == contentType {
				return codec, nil
			}
		}
	}
	return nil, NewResponseError(http.StatusNotFound, "codec not found")
}

// FindByFileExtensionOrContentType finds a codec by file extension or content type
func (c CodecRegistry) FindByFileExtensionOrContentType(ext string, contentType string) (*Codec, error) {
	if codec, err := c.FindByFileExtension(ext); err == nil {
		return codec, nil
	}
	return c.FindByContentType(contentType)
}

// CODEC IMPLEMENTATIONS
var DefaultCodecs = []*Codec{
	{
		Id:             "json",
		FileExtensions: []string{".json"},
		ContentTypes:   []string{"application/json"},
		OnDeserialize: func(r io.Reader) (map[string]any, error) {
			var data map[string]any
			err := json.NewDecoder(r).Decode(&data)
			if err != nil {
				return nil, err
			}
			return data, nil
		},
		OnSerialize: func(record *Record) ([]byte, error) {
			return json.Marshal(record.Data)
		},
	},
	{
		Id:             "markdown",
		FileExtensions: []string{".md", ".markdown"},
		ContentTypes:   []string{"text/markdown"},
		OnDeserialize: func(r io.Reader) (map[string]any, error) {
			content, err := io.ReadAll(r)
			if err != nil {
				return nil, err
			}
			return map[string]any{
				"content": string(content),
			}, nil
		},
		OnSerialize: func(record *Record) ([]byte, error) {
			return []byte(record.Data["content"].(string)), nil
		},
	},
}
