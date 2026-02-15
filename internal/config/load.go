package config

import (
	"fmt"
	"os"

	"sigs.k8s.io/yaml"
)

func LoadApp(path string) (*App, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var a App
	if err := yaml.Unmarshal(b, &a); err != nil {
		return nil, err
	}

	// minimal validation + defaults
	if a.APIVersion == "" {
		a.APIVersion = "gpp.dev/v1alpha1"
	}
	if a.Kind == "" {
		a.Kind = "App"
	}

	if a.Metadata.Name == "" {
		return nil, fmt.Errorf("metadata.name is required")
	}
	if a.Metadata.Namespace == "" {
		return nil, fmt.Errorf("metadata.namespace is required")
	}
	if a.Spec.Image == "" {
		return nil, fmt.Errorf("spec.image is required")
	}
	if a.Spec.Port == 0 {
		a.Spec.Port = 8080
	}
	return &a, nil
}
