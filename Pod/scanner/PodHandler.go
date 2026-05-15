package pod

import "fmt"

//BootUpPod is the entry point for the scanner which loads registry and passes it to scan finction
func BootUpPod(regPath string) {
	reg, err := loadRegistry(regPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Registry: %s | Version: %v | Platform: %s", reg.Schema.Name, reg.Schema.Version, reg.Platform)

	//If fine then lets start the Scanning
	scanData := ScanAnamolies(reg)
	archivePath, err := WriteScanArchive(scanData)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("\nScan JSON written to %s\n", archivePath)

}
