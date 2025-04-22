package helper

import (
	"os"

	"helm.sh/helm/v3/pkg/chartutil"
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
