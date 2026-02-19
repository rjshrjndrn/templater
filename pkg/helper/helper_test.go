package helper

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestMergeYaml_DeepMergeByDefault(t *testing.T) {
	dst := map[string]any{
		"class": map[string]any{
			"name": "suresh",
		},
	}
	src := map[string]any{
		"class": map[string]any{
			"name":  "rajesh",
			"grade": 10,
		},
	}
	pruneReplaced(dst, src)
	MergeYaml(dst, src)
	class := dst["class"].(map[string]any)
	if class["name"] != "suresh" {
		t.Errorf("expected name=suresh, got %v", class["name"])
	}
	if class["grade"] != 10 {
		t.Errorf("expected grade=10, got %v", class["grade"])
	}
}

func TestMergeYaml_ReplaceBlocksMerge(t *testing.T) {
	dst := map[string]any{
		"class": map[string]any{
			"__replace": true,
			"name":      "suresh",
		},
	}
	src := map[string]any{
		"class": map[string]any{
			"name":  "rajesh",
			"grade": 10,
		},
	}
	pruneReplaced(dst, src)
	MergeYaml(dst, src)
	class := dst["class"].(map[string]any)
	if class["name"] != "suresh" {
		t.Errorf("expected name=suresh, got %v", class["name"])
	}
	if _, ok := class["grade"]; ok {
		t.Error("grade should not be present when __replace is used")
	}
}

func TestMergeYaml_ThreeFiles_ReplaceInMiddle(t *testing.T) {
	// Simulates: -f val1.yaml -f val2.yaml -f val3.yaml
	// val2 has __replace on class, val3 merges normally
	// Iteration order: val3 (i=2), val2 (i=1), val1 (i=0)

	val1 := map[string]any{
		"class": map[string]any{
			"name":  "rajesh",
			"grade": 10,
		},
	}
	val2 := map[string]any{
		"class": map[string]any{
			"__replace": true,
			"name":      "suresh",
		},
	}
	val3 := map[string]any{
		"class": map[string]any{
			"section": "A",
		},
	}

	// Start with val3 (highest precedence)
	data := val3
	// Merge val2 (val3 takes precedence)
	pruneReplaced(data, val2)
	MergeYaml(data, val2)
	// Merge val1 (data takes precedence, __replace blocks val1)
	pruneReplaced(data, val1)
	MergeYaml(data, val1)

	StripReplaceAnnotations(data)

	class := data["class"].(map[string]any)
	if class["section"] != "A" {
		t.Errorf("expected section=A, got %v", class["section"])
	}
	if class["name"] != "suresh" {
		t.Errorf("expected name=suresh (from val2), got %v", class["name"])
	}
	if _, ok := class["grade"]; ok {
		t.Error("grade from val1 should not be present")
	}
	if _, ok := class["__replace"]; ok {
		t.Error("__replace should be stripped")
	}
}

func TestMergeYaml_NestedReplace(t *testing.T) {
	dst := map[string]any{
		"app": map[string]any{
			"config": map[string]any{
				"__replace": true,
				"port":      8080,
			},
		},
	}
	src := map[string]any{
		"app": map[string]any{
			"config": map[string]any{
				"port": 3000,
				"host": "localhost",
			},
			"version": "1.0",
		},
	}
	pruneReplaced(dst, src)
	MergeYaml(dst, src)
	StripReplaceAnnotations(dst)

	app := dst["app"].(map[string]any)
	config := app["config"].(map[string]any)

	if config["port"] != 8080 {
		t.Errorf("expected port=8080, got %v", config["port"])
	}
	if _, ok := config["host"]; ok {
		t.Error("host should not be present when __replace is used")
	}
	// version should still merge at app level
	if app["version"] != "1.0" {
		t.Errorf("expected version=1.0, got %v", app["version"])
	}
}

func TestStripReplaceAnnotations(t *testing.T) {
	m := map[string]any{
		"__replace": true,
		"a":         "b",
		"nested": map[string]any{
			"__replace": true,
			"c":         "d",
		},
	}
	StripReplaceAnnotations(m)

	if _, ok := m["__replace"]; ok {
		t.Error("top-level __replace should be stripped")
	}
	nested := m["nested"].(map[string]any)
	if _, ok := nested["__replace"]; ok {
		t.Error("nested __replace should be stripped")
	}
	if m["a"] != "b" || nested["c"] != "d" {
		t.Error("non-annotation keys should be preserved")
	}
}

func TestParseYAMLValues_ReplaceIntegration(t *testing.T) {
	dir := t.TempDir()

	val1 := filepath.Join(dir, "val1.yaml")
	val2 := filepath.Join(dir, "val2.yaml")

	os.WriteFile(val1, []byte("class:\n  name: rajesh\n  grade: 10\n"), 0644)
	os.WriteFile(val2, []byte("class:\n  __replace: true\n  name: suresh\n"), 0644)

	result, err := ParseYAMLValues([]string{val1, val2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	values := result["Values"].(map[string]any)
	class := values["class"].(map[string]any)

	expected := map[string]any{"name": "suresh"}
	if !reflect.DeepEqual(class, expected) {
		t.Errorf("expected %v, got %v", expected, class)
	}
}

func TestParseYAMLValues_ThreeFiles_ReplaceInMiddle(t *testing.T) {
	dir := t.TempDir()

	val1 := filepath.Join(dir, "val1.yaml")
	val2 := filepath.Join(dir, "val2.yaml")
	val3 := filepath.Join(dir, "val3.yaml")

	os.WriteFile(val1, []byte("class:\n  name: rajesh\n  grade: 10\n"), 0644)
	os.WriteFile(val2, []byte("class:\n  __replace: true\n  name: suresh\n"), 0644)
	os.WriteFile(val3, []byte("class:\n  section: A\n"), 0644)

	result, err := ParseYAMLValues([]string{val1, val2, val3})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	values := result["Values"].(map[string]any)
	class := values["class"].(map[string]any)

	if class["section"] != "A" {
		t.Errorf("expected section=A from val3, got %v", class["section"])
	}
	if class["name"] != "suresh" {
		t.Errorf("expected name=suresh from val2, got %v", class["name"])
	}
	if _, ok := class["grade"]; ok {
		t.Error("grade from val1 should not be present")
	}
}
