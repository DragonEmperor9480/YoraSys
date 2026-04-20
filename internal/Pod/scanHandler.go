package pod

import (
	"errors"
	"fmt"
	"os"

	schematics "github.com/DragonEmperor9480/yorasys/internal/Schematics"
	"gopkg.in/yaml.v3"
)

func ScanAnamolies(registryPath string) {
	reg, err := loadRegistry(registryPath)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Registry:", reg)

}

func loadRegistry(path string) (schematics.Registry, error) {
	val, err := os.ReadFile(path)
	if err != nil {
		return schematics.Registry{}, errors.New("Failed to Read the Registry")
	}

	var reg schematics.Registry

	if err := yaml.Unmarshal(val, &reg); err != nil {
		return schematics.Registry{}, errors.New("Failed to Parse yaml")

	}

	if reg.Platform == "" {
		return schematics.Registry{}, errors.New("Invalid OS")
	}
	if len(reg.Caches) == 0 {
		return schematics.Registry{}, errors.New("No cache data found")

	}

	return reg, nil

}
