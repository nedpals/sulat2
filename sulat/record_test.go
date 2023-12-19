package sulat

import "testing"

func TestParseField(t *testing.T) {
	test := map[string]any{
		"a": map[string]any{
			"b": []any{1, map[string]any{"c": []any{2}}},
		},
	}

	fp := parseField("a.b.1.c.0", test)
	value := fp.get()

	if value != 2 {
		t.Errorf("Expected 2, got %v", value)
	}

	// non existent key
	fp2 := parseField("a.d", test)
	value2 := fp2.get()

	if value2 != nil {
		t.Errorf("Expected non nil value")
	}
}
