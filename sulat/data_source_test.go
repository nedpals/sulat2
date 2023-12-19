package sulat

import (
	"io"
	"testing"
	"testing/fstest"

	"github.com/nedpals/sulatcms/sulat/query"
	"github.com/spf13/afero"
)

func TestFileDataSourceProvider(t *testing.T) {
	inst, err := NewInstance("")
	if err != nil {
		t.Fatal(err)
	}

	if err := inst.RegisterCodecs(DefaultCodecs...); err != nil {
		t.Fatal(err)
	}

	testFs := afero.NewCopyOnWriteFs(afero.FromIOFS{
		FS: fstest.MapFS{
			"project/data/socials.json": {
				Data: []byte(`{"data": [{"id": "facebook", "name": "Facebook"}]}`),
			},
			"project/posts/hello-world.md": {
				Data: []byte(`---
title: Hello World
---
Hello!
`),
			},
			"sulat.toml": {
				Data: []byte(`
[collections]
data = "project/data/*.json"
posts = "project/posts/*.md"
`),
			},
		},
	}, afero.NewMemMapFs())

	var provider DataSourceProvider = &FileDataSourceProvider{FS: testFs}

	t.Run("Simple", func(t *testing.T) {
		dataSource := inst.NewDataSource("sample", "Sample", provider, map[string]any{
			"root": "project",
			"collections": map[string]string{
				"data":  "data/*.json",
				"posts": "posts/*.md",
			},
		})

		// Sample find
		records, err := dataSource.Find("posts", query.Eq("id", "hello-world.md"), nil)
		if err != nil {
			t.Fatal(err)
		} else if len(records) != 1 {
			t.Fatalf("Expected 1 record, got %d", len(records))
		}

		// ...how about for the data collection?
		records, err = dataSource.Find("data", query.Eq("data.0.id", "facebook"), nil)
		if err != nil {
			t.Fatal(err)
		} else if len(records) != 1 {
			t.Fatalf("Expected 1 record, got %d", len(records))
		} else if records[0].Id != "socials.json" {
			t.Fatalf("Expected id to be 'socials.json', got %s", records[0].Id)
		}

		// Sample get
		record, err := dataSource.Get("posts", "hello-world.md", nil)
		if err != nil {
			t.Fatal(err)
		}

		if record.Id != "hello-world.md" {
			t.Fatalf("Expected id to be 'hello-world.md', got %s", record.Id)
		}

		// Sample insert
		err = dataSource.Insert("posts", &Record{
			Id:   "foo-bar.md",
			Data: map[string]any{"content": "## Foo Bar!"},
		}, nil)
		if err != nil {
			t.Fatal(err)
		}

		records, err = dataSource.Find("posts", query.Eq("id", "foo-bar.md"), nil)
		if err != nil {
			t.Fatal(err)
		}

		if len(records) != 1 {
			t.Fatalf("Expected 1 record, got %d", len(records))
		}

		if _, err := testFs.Stat("project/posts/foo-bar.md"); err != nil {
			t.Fatal(err)
		}

		// Sample update
		err = dataSource.Update("posts", &Record{
			Id:   "foo-bar.md",
			Data: map[string]any{"content": "## Bar baz!"},
		}, nil)
		if err != nil {
			t.Fatal(err)
		}

		record, err = dataSource.Get("posts", "foo-bar.md", nil)
		if err != nil {
			t.Fatal(err)
		}

		if record.Data["content"] != "## Bar baz!" {
			t.Fatalf("Expected content to be '## Bar baz!', got %s", record.Data["content"])
		}

		// check if file has been updated in testfs
		if _, err := testFs.Stat("project/posts/foo-bar.md"); err != nil {
			t.Fatal(err)
		}

		// open file to see if the contents have been updated
		file, err := testFs.Open("project/posts/foo-bar.md")
		if err != nil {
			t.Fatal(err)
		}

		content, err := io.ReadAll(file)
		if err != nil {
			t.Fatal(err)
		}

		if string(content) != "## Bar baz!" {
			t.Fatalf("Expected content to be '## Bar baz!', got %s", string(content))
		}

		// Sample delete
		err = dataSource.Delete("posts", query.Eq("id", "bar-baz.md"), nil)
		if err != nil {
			t.Fatal(err)
		}

		// ... check if file has been deleted
		if _, err := testFs.Stat("project/posts/bar-baz.md"); err == nil {
			t.Fatal("Expected file to be deleted, but it wasn't")
		}

		// ... check if the post has been deleted successfully
		records, err = dataSource.Find("posts", query.Eq("id", "bar-baz.md"), nil)
		if err != nil && err.Error() != "no records found" {
			t.Fatal(err)
		}

		if len(records) != 0 {
			t.Fatalf("Expected 0 record, got %d", len(records))
		}
	})

	t.Run("With config file", func(t *testing.T) {
		dataSource := inst.NewDataSource("sample", "Sample", provider, map[string]any{
			"config_path": "sulat.toml",
		})

		// Sample find
		records, err := dataSource.Find("posts", query.Eq("id", "hello-world.md"), nil)
		if err != nil {
			t.Fatal(err)
		} else if len(records) != 1 {
			t.Fatalf("Expected 1 record, got %d", len(records))
		}
	})
}
