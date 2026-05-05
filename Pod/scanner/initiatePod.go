package pod

import (
	"errors"
	"os"

	schematics "github.com/DragonEmperor9480/yorasys/Pod/Schematics"
	"gopkg.in/yaml.v3"
)

func loadRegistry(path string) (schematics.Registry, error) {
	val, err := os.ReadFile(path)
	if err != nil {
		return schematics.Registry{}, errors.New("Failed to Read the Registry")
	}

	var reg schematics.Registry

	if err := yaml.Unmarshal(val, &reg); err != nil {
		return reg, errors.New("Failed to Parse yaml")

	}

	if reg.Platform == "" {
		return reg, errors.New("Invalid OS")
	}
	if len(reg.Caches) == 0 {
		return reg, errors.New("No cache data found")

	}

	return reg, nil

}
