package helper

import (
	"fmt"
	"os"

	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/strvals"
	"sigs.k8s.io/yaml"
)

// parseYAMLValues reads a YAML file and returns the parsed data.
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
		MergeYaml(data, values)
	}

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
