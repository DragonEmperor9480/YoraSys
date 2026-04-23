package pod

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	schematics "github.com/DragonEmperor9480/yorasys/internal/Schematics"
	"gopkg.in/yaml.v3"
)

func ScanAnamolies(registryPath string) {
	reg, err := loadRegistry(registryPath)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Registry: %s | Version: %v | Platform: %s", reg.Schema.Name, reg.Schema.Version, reg.Platform)

	for _, valCache := range reg.Caches {
		cachePresent := false
		fmt.Printf("\nCache: %s (ID: %d)\n", valCache.Name, valCache.ID)

		for _, cachePath := range valCache.Paths {
			subPaths, err := handleFullPath(cachePath)
			if err != nil {
				fmt.Printf("Wrong Yaml data on %v, err: %v\n", cachePath, err)
			}

			if len(subPaths) == 0 {
				fmt.Printf("am never gonna execute but justtt lets see", cachePath)
			}
			exists, _, err := checkPath(cachePath)
			if err != nil {
				fmt.Printf("program.exe is meow meow %v\n", err)
				continue
			}
			if exists {
				cachePresent = true
				fmt.Printf("Found something\n", cachePath)
			} else {
				fmt.Printf("Meh didnt find a thing\n", cachePath)
			}
		}
		fmt.Printf("hmmmmm %v\n", cachePresent)
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

func checkPath(path string) (exists bool, isDir bool, err error) {
	info, err := os.Stat(path)
	if err == nil {
		return true, info.IsDir(), nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return false, false, nil
	}
	return false, false, err
}

func handleFullPath(path string) ([]string, error) {
	if !strings.Contains(path, "*") {
		return []string{path}, nil
	}

	found, err := filepath.Glob(path)
	if err != nil {
		return nil, err
	}

	return found, nil
}
