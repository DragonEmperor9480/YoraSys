package main

import (
	pod "github.com/DragonEmperor9480/yorasys/Pod/scanner"
)

const defaultRegistryPath = "registry/scanData_windows.yaml"

func main() {
	pod.BootUpScanPod(defaultRegistryPath)

}
