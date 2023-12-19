package query

import (
	"encoding/json"
	"testing"

	"github.com/go-test/deep"
)

func TestParser(t *testing.T) {
	// sample inputs:
	// - and(eq(id 1),eq(name "John Doe"),{limit:10, offset:0, order:["id","desc"]})
	// - or(eq(id 1),eq(name "John Doe"))
	// - eq(id 1)
	parser := NewParser()
	parser.IsTest = true

	t.Run("Test parsing simple query", func(t *testing.T) {
		query := "eq(id 1)"
		expected := Query{
			Operator: "eq",
			Field:    "id",
			Value:    json.Number("1"),
			Options:  nil,
		}

		result, err := parser.Parse(query)
		if err != nil {
			t.Errorf("Error parsing query: %v", err)
		}

		if diff := deep.Equal(result.Options, expected.Options); diff != nil {
			t.Error(diff)
		}
	})

	t.Run("Test parsing AND query", func(t *testing.T) {
		query := "and(eq(id 1), eq(name \"John Doe\"),{limit:10, offset:0, order:[\"id\",\"desc\"]})"
		expected := &Query{
			Operator: "and",
			Value: []*Query{
				{
					Operator: "eq",
					Field:    "id",
					Value:    json.Number("1"),
				},
				{
					Operator: "eq",
					Field:    "name",
					Value:    "John Doe",
				},
			},
			Options: map[string]any{
				"limit":  json.Number("10"),
				"offset": json.Number("0"),
				"order":  []any{"id", "desc"},
			},
		}

		result, err := parser.Parse(query)
		if err != nil {
			t.Errorf("Error parsing query: %v", err)
		}

		if diff := deep.Equal(result, expected); diff != nil {
			t.Error(diff)
		}
	})

	t.Run("Test parsing OR query", func(t *testing.T) {
		query := "or(eq(id 1),eq(name \"John Doe\"))"
		expected := &Query{
			Operator: "or",
			Value: []*Query{
				{
					Operator: "eq",
					Field:    "id",
					Value:    json.Number("1"),
				},
				{
					Operator: "eq",
					Field:    "name",
					Value:    "John Doe",
				},
			},
		}

		result, err := parser.Parse(query)
		if err != nil {
			t.Errorf("Error parsing query: %v", err)
		}

		if diff := deep.Equal(result, expected); diff != nil {
			t.Error(diff)
		}
	})

	t.Run("Test simple query with options", func(t *testing.T) {
		query := "eq(id 1, {limit:10, offset:0, order:[\"id\",\"desc\"]})"
		expected := &Query{
			Operator: "eq",
			Field:    "id",
			Value:    json.Number("1"),
			Options: map[string]any{
				"limit":  json.Number("10"),
				"offset": json.Number("0"),
				"order":  []any{"id", "desc"},
			},
		}

		result, err := parser.Parse(query)
		if err != nil {
			t.Errorf("Error parsing query: %v", err)
		}

		if diff := deep.Equal(result, expected); diff != nil {
			t.Error(diff)
		}
	})

	t.Run("Test query with nested options", func(t *testing.T) {
		query := "eq(id 1, {limit:10, offset:0, order:[\"id\",\"desc\"], custom: {test: 123, deep: {hello:\"world\"}}})"
		expected := &Query{
			Operator: "eq",
			Field:    "id",
			Value:    json.Number("1"),
			Options: map[string]any{
				"limit":  json.Number("10"),
				"offset": json.Number("0"),
				"order":  []any{"id", "desc"},
				"custom": map[string]any{
					"test": json.Number("123"),
					"deep": map[string]any{
						"hello": "world",
					},
				},
			},
		}

		result, err := parser.Parse(query)
		if err != nil {
			t.Errorf("Error parsing query: %v", err)
		}

		if diff := deep.Equal(result, expected); diff != nil {
			t.Error(diff)
		}
	})

	t.Run("Test query with dot notation field", func(t *testing.T) {
		query := "eq(a.b.1.c.0 1)"
		expected := &Query{
			Operator: "eq",
			Field:    "a.b.1.c.0",
			Value:    json.Number("1"),
		}

		result, err := parser.Parse(query)
		if err != nil {
			t.Errorf("Error parsing query: %v", err)
		}

		if diff := deep.Equal(result, expected); diff != nil {
			t.Error(diff)
		}
	})

	t.Run("Test deeply nested query w/o options", func(t *testing.T) {
		query := "and(eq(id 1),eq(name \"John Doe\"),or(eq(id 1),eq(name \"John Doe\")))"
		expected := &Query{
			Operator: "and",
			Value: []*Query{
				{
					Operator: "eq",
					Field:    "id",
					Value:    json.Number("1"),
				},
				{
					Operator: "eq",
					Field:    "name",
					Value:    "John Doe",
				},
				{
					Operator: "or",
					Value: []*Query{
						{
							Operator: "eq",
							Field:    "id",
							Value:    json.Number("1"),
						},
						{
							Operator: "eq",
							Field:    "name",
							Value:    "John Doe",
						},
					},
				},
			},
		}

		result, err := parser.Parse(query)
		if err != nil {
			t.Errorf("Error parsing query: %v", err)
		}

		if diff := deep.Equal(result, expected); diff != nil {
			t.Error(diff)
		}
	})
}
