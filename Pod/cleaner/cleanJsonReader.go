package cleaner

import (
	"encoding/json"
	"fmt"
	"os"

	schematics "github.com/DragonEmperor9480/yorasys/Pod/Schematics"
)

func LoadCleanJson(path string) (schematics.CleanSelection, error) {
	val, err := os.ReadFile(path)
	if err != nil {
		return schematics.CleanSelection{}, fmt.Errorf("failed to read clean json: %w", err)
	}

	var cleanData schematics.CleanSelection
	if err := json.Unmarshal(val, &cleanData); err != nil {
		return schematics.CleanSelection{}, fmt.Errorf("failed to parse clean json: %w", err)
	}

	return cleanData, nil
}
