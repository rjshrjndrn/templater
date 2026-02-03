package utils

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"sigs.k8s.io/yaml"
)

func ToYAMLFunc(obj interface{}) (string, error) {
	data, err := yaml.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// MakeTplFunc creates a tpl function that can evaluate strings as templates
func MakeTplFunc() func(string, any) (string, error) {
	var tplFunc func(string, any) (string, error)
	tplFunc = func(tplString string, data any) (string, error) {
		nestedFm := sprig.FuncMap()
		nestedFm["toYaml"] = ToYAMLFunc
		nestedFm["tpl"] = tplFunc
		t, err := template.New("tpl").Funcs(nestedFm).Parse(tplString)
		if err != nil {
			return "", err
		}
		var buf bytes.Buffer
		if err := t.Execute(&buf, data); err != nil {
			return "", err
		}
		return buf.String(), nil
	}
	return tplFunc
}

// MakeIncludeFunc creates an include function that can load and execute external template files
func MakeIncludeFunc(baseDir string, tplFunc func(string, any) (string, error)) func(string, any) (string, error) {
	var makeIncludeFunc func(string) func(string, any) (string, error)
	makeIncludeFunc = func(currentBaseDir string) func(string, any) (string, error) {
		return func(filePath string, data any) (string, error) {
			// Resolve path relative to the current template's directory
			var resolvedPath string
			if filepath.IsAbs(filePath) {
				resolvedPath = filePath
			} else {
				resolvedPath = filepath.Join(currentBaseDir, filePath)
			}

			// Read the template file
			includeContent, err := os.ReadFile(resolvedPath)
			if err != nil {
				return "", fmt.Errorf("failed to read include file %s: %w", resolvedPath, err)
			}

			// Get the directory of the included file for nested includes
			includeBaseDir := filepath.Dir(resolvedPath)

			// Create function map for the included template
			includeFm := sprig.FuncMap()
			includeFm["toYaml"] = ToYAMLFunc
			includeFm["tpl"] = tplFunc
			includeFm["include"] = makeIncludeFunc(includeBaseDir)

			// Parse and execute the included template
			t, err := template.New(filepath.Base(resolvedPath)).Funcs(includeFm).Parse(string(includeContent))
			if err != nil {
				return "", fmt.Errorf("failed to parse include file %s: %w", resolvedPath, err)
			}

			var buf bytes.Buffer
			if err := t.Execute(&buf, data); err != nil {
				return "", fmt.Errorf("failed to execute include file %s: %w", resolvedPath, err)
			}
			return buf.String(), nil
		}
	}
	return makeIncludeFunc(baseDir)
}
