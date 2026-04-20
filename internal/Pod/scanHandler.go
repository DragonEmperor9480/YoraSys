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

	fmt.Printf("Registry: %s | Version: %v | Platform: %s", reg.Schema.Name, reg.Schema.Version, reg.Platform)

	for _, cache := range reg.Caches {
		fmt.Printf("ID: %d\n", cache.ID)
		fmt.Printf("Name: %s\n", cache.Name)
		fmt.Printf("Category: %s\n", cache.Category)
		fmt.Printf("Description: %s\n", cache.Description)

		if len(cache.Paths) == 0 {
			fmt.Println("Paths: none")
			continue
		}

		fmt.Println("Paths:")
		for j, p := range cache.Paths {
			fmt.Printf("  %d. %s\n", j+1, p)
		}

		fmt.Println()
	}

}

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
