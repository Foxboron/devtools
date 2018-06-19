package main

import (
	"flag"
	"os"

	"github.com/foxboron/devtools/backend"
	"github.com/foxboron/devtools/backend/fs/overlay"
	"github.com/foxboron/devtools/utils"
)

func main() {

	flag.Parse()
	if flag.NArg() < 2 {
		utils.Error("You must specify a directory.")
		os.Exit(1)
	}

	WorkingDir := flag.Args()[0]

	var backendFs backend.BackendFilesystem
	backendFs, err := backend.GetBackendFromContainer(WorkingDir)
	if err != nil {
		utils.Error("Not a container!")
		os.Exit(1)
	}

	var backendInit backend.Backend
	switch backendFs {
	case backend.Overlay:
		backendInit = overlay.NewOverlay(WorkingDir)
	default:
		utils.Error("Invalid filesystem")
		os.Exit(1)
	}

	if err := backendInit.Destroy(); err != nil {
		utils.Error("Failed to destroy backend")
		os.Exit(1)
	}
}
