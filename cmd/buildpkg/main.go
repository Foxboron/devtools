package main

import (
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/foxboron/devtools/backend"
	"github.com/foxboron/devtools/backend/fs/overlay"
	"github.com/foxboron/devtools/bootstrap"
	"github.com/foxboron/devtools/bootstrap/archiso"
	"github.com/foxboron/devtools/bootstrap/pacstrap"
	"github.com/foxboron/devtools/builder"
	"github.com/foxboron/devtools/container/nspawn"
	"github.com/foxboron/devtools/utils"
)

// Optional values
const (
	BuildPath    string = "/var/lib/buildpkg"
	Architecture        = "x86_64"
	Repository          = "extra"
	Backend             = "overlay"
	Bootstrap           = "archiso"
	PacmanConf          = "/etc/pacman.conf"
	MakepkgConf         = "/etc/makepkg.conf"
)

func main() {
	if usr, _ := user.Current(); usr.Uid != "0" {
		utils.Error("Need to be run as sudo")
		os.Exit(1)
	}

	var backendInit backend.Backend
	switch backend.GetBackend(Backend) {
	case backend.Overlay:
		backendInit = overlay.NewOverlay(BuildPath)
	default:
		utils.Error("Invalid filesystem")
		os.Exit(1)
	}

	var bootstrapInit bootstrap.Bootstrap
	switch bootstrap.GetBootstrap(Bootstrap) {
	case bootstrap.Archiso:
		bootstrapInit = archiso.NewArchiso(PacmanConf)
	case bootstrap.Pacstrap:
		bootstrapInit = pacstrap.NewPacstrap(PacmanConf)
	}

	build := &builder.Builder{
		Path:         BuildPath,
		Backend:      backendInit,
		Bootstrap:    bootstrapInit,
		Container:    nspawn.NewNspawn(""),
		Repository:   Repository,
		Architecture: Architecture,
		MakepkgConf:  MakepkgConf,
		PacmanConf:   PacmanConf,
	}

	err := build.Init()
	if err != nil {
		log.Fatal(err)
	}
	err = build.Update()
	if err != nil {
		log.Fatal(err)
	}

	containerName := os.Getenv("SUDO_USER")
	if containerName == "" {
		log.Fatal("Couldn't get USER name!")
	}

	err = build.Fork(containerName)
	if err != nil {
		build.Destroy(containerName)
		log.Fatal(err)
	}
	utils.Msg(fmt.Sprintf("Synchronizing chroot copy [%s] -> [%s]", "root", containerName))
	// build.SetEnv(*Environment)
	err = build.Build()
	if err != nil {
		log.Fatal(err)
	}
	utils.Msg(fmt.Sprintf("Deleting chroot copy [%s]", containerName))
	build.Destroy(containerName)
	os.Exit(0)
	return
}
