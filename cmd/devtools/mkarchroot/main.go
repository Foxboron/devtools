package main

import (
	"flag"
	"log"
	"os"

	"github.com/foxboron/devtools/backend"
	"github.com/foxboron/devtools/backend/fs/overlay"
	"github.com/foxboron/devtools/bootstrap"
	"github.com/foxboron/devtools/bootstrap/archiso"
	"github.com/foxboron/devtools/bootstrap/pacstrap"
	"github.com/foxboron/devtools/builder"
	"github.com/foxboron/devtools/container/nspawn"
	"github.com/foxboron/devtools/utils"
)

type listString []string

func (i *listString) String() string {
	return "my string representation"
}

func (i *listString) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var (
	CopyFiles     listString
	BootstrapType = flag.String("b", "archiso", "Bootstrap method. One of:")
	PacmanCache   = flag.String("c", "/var/cache/pacman/pkg", "Set pacman cache")
	BackendType   = flag.String("t", "overlay", "Backend type. One of:")
	Help          = flag.Bool("h", false, "This message")
	SetArch       = flag.Bool("s", true, "Do not run setarch")
	PacmanConf    = flag.String("C", "/etc/pacman.conf", "Location of a pacman config file")
	MakepkgConf   = flag.String("M", "/etc/makepkg.conf", "Location of a makepkg config file")
)

func main() {
	flag.Var(&CopyFiles, "f", "Copy file from the host to the chroot")

	flag.Usage = func() {
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 2 {
		panic("You must specify a directory and one or more packages.")
	}

	WorkingDir := flag.Args()[0]
	// Packages := flag.Args()[1:]

	var backendInit backend.Backend
	switch backend.GetBackend(*BackendType) {
	case backend.Overlay:
		backendInit = overlay.NewOverlay(WorkingDir)
	default:
		utils.Error("Invalid filesystem")
		os.Exit(1)
	}
	var bootstrapInit bootstrap.Bootstrap
	switch bootstrap.GetBootstrap(*BootstrapType) {
	case bootstrap.Archiso:
		bootstrapInit = archiso.NewArchiso(string(*PacmanConf))
	case bootstrap.Pacstrap:
		bootstrapInit = pacstrap.NewPacstrap(string(*PacmanConf))
	}

	build := &builder.Builder{
		Path:        WorkingDir,
		Backend:     backendInit,
		Bootstrap:   bootstrapInit,
		Container:   nspawn.NewNspawn(WorkingDir),
		MakepkgConf: string(*MakepkgConf),
		PacmanConf:  string(*PacmanConf),
	}
	err := build.Init()
	if err != nil {
		log.Fatal(err)
	}
	utils.Msgf("Created arch chroot at %s", WorkingDir)
	os.Exit(0)
}
