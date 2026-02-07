package helper

import (
	"os"
	"strconv"
	"strings"

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

// ParseSetValues parses --set flags like Helm (e.g., "apple.eat=true")
// and returns a map that can be merged with other values.
func ParseSetValues(setValues []string) (map[string]any, error) {
	result := make(map[string]any)

	for _, setValue := range setValues {
		// Split on first '=' to get key and value
		parts := strings.SplitN(setValue, "=", 2)
		if len(parts) != 2 {
			continue // Skip malformed entries
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		// Parse the value to appropriate type
		parsedVal := parseValue(val)

		// Set the nested value using dot notation
		setNestedValue(result, key, parsedVal)
	}

	return result, nil
}

// parseValue attempts to parse a string value into its appropriate type
func parseValue(val string) any {
	// Check for boolean
	if val == "true" {
		return true
	}
	if val == "false" {
		return false
	}

	// Check for null/nil
	if val == "null" || val == "nil" {
		return nil
	}

	// Check for integer
	if i, err := strconv.ParseInt(val, 10, 64); err == nil {
		return i
	}

	// Check for float
	if f, err := strconv.ParseFloat(val, 64); err == nil {
		return f
	}

	// Return as string
	return val
}

// setNestedValue sets a value in a nested map using dot notation key
func setNestedValue(m map[string]any, key string, value any) {
	keys := strings.Split(key, ".")
	current := m

	for i, k := range keys {
		if i == len(keys)-1 {
			// Last key, set the value
			current[k] = value
		} else {
			// Create nested map if it doesn't exist
			if _, exists := current[k]; !exists {
				current[k] = make(map[string]any)
			}
			// Navigate deeper (create new map if existing value isn't a map)
			if nested, ok := current[k].(map[string]any); ok {
				current = nested
			} else {
				newMap := make(map[string]any)
				current[k] = newMap
				current = newMap
			}
		}
	}
}
