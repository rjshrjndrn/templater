package utils

import "sigs.k8s.io/yaml"

func ToYAMLFunc(obj interface{}) (string, error) {
	data, err := yaml.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
