package helper

import (
	"fmt"
	"os"

	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/strvals"
	"sigs.k8s.io/yaml"
)

const replaceAnnotation = "__replace"

// ParseYAMLValues reads YAML files and returns the merged data.
// Later files take precedence over earlier ones.
// If a map contains "__replace: true", it replaces the corresponding map
// from earlier files entirely instead of deep merging.
func ParseYAMLValues(valuesPath []string) (map[string]any, error) {
	if len(valuesPath) == 0 {
		return map[string]any{}, nil // No values file provided
	}

	data := make(map[string]any)
	for i := len(valuesPath) - 1; i >= 0; i-- {
		path := valuesPath[i]
		values := make(map[string]any)
		d, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		err = yaml.Unmarshal(d, &values)
		if err != nil {
			return nil, err
		}
		pruneReplaced(data, values)
		MergeYaml(data, values)
	}

	StripReplaceAnnotations(data)

	// Support helm structure
	// {{ .key }} -> {{ .Values.key }}
	return map[string]any{"Values": data}, nil
}

// Merge map2 to map1:
//
// map1 takes precedence
func MergeYaml(map1, map2 map[string]any) map[string]any {
	return chartutil.CoalesceTables(map1, map2)
}

// pruneReplaced walks dst and src in parallel. For any key where dst holds
// a map with __replace: true, the corresponding key is deleted from src
// so that CoalesceTables will not deep merge it.
func pruneReplaced(dst, src map[string]any) {
	for key, dstVal := range dst {
		srcVal, exists := src[key]
		if !exists {
			continue
		}
		dstMap, dstIsMap := dstVal.(map[string]any)
		if !dstIsMap {
			continue
		}
		if isReplace(dstMap) {
			delete(src, key)
			continue
		}
		// Recurse into nested maps
		if srcMap, srcIsMap := srcVal.(map[string]any); srcIsMap {
			pruneReplaced(dstMap, srcMap)
		}
	}
}

// isReplace checks if a map has the __replace annotation set to true.
func isReplace(m map[string]any) bool {
	v, ok := m[replaceAnnotation]
	if !ok {
		return false
	}
	if b, ok := v.(bool); ok {
		return b
	}
	return false
}

// StripReplaceAnnotations recursively removes all "__replace" keys from the map
// so they don't leak into the template context.
func StripReplaceAnnotations(m map[string]any) {
	delete(m, replaceAnnotation)
	for _, v := range m {
		if subMap, ok := v.(map[string]any); ok {
			StripReplaceAnnotations(subMap)
		}
	}
}

// ParseSetValues parses --set flags using Helm's strvals parser.
// Supports full Helm --set syntax including dot notation, arrays, and escaping.
func ParseSetValues(setValues []string, dest map[string]any) error {
	for _, value := range setValues {
		if err := strvals.ParseInto(value, dest); err != nil {
			return fmt.Errorf("failed parsing --set data: %w", err)
		}
	}
	return nil
}
