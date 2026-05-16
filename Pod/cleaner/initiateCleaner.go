package cleaner

import "fmt"

func BootUpCleaner(cleanPath string) {
	cleanData, err := LoadCleanJson(cleanPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Clean JSON loaded: %s | selected_roots: %d\n", cleanPath, len(cleanData.SelectedPaths))
}
