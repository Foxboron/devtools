package builder

import (
	"fmt"
	"os"
	"path"

	"github.com/foxboron/devtools/backend"
	"github.com/foxboron/devtools/bootstrap"
	"github.com/foxboron/devtools/container"
	"github.com/foxboron/devtools/makepkg"
	"github.com/foxboron/devtools/utils"
)

const (
	makepkgArgs    = "syncdeps --noconfirm --log --holdver --skipinteg"
	makepkgCommand = "sudo --preserve-env=SOURCE_DATE_EPOCH -iu builduser bash -c 'cd /startdir; makepkg \"$@\"' -bash "
)

var (
	// Fixed paths
	hostGnupgPath = path.Join("/etc", "pacman.d", "gnupg")
)

type Builder struct {
	Bootstrap     bootstrap.Bootstrap
	Backend       backend.Backend
	Container     container.Container
	Path          string
	ContainerPath string
	Repository    string
	Architecture  string
	MakepkgConf   string
	PacmanConf    string
}

func (b *Builder) Build() (map[string]map[string]string, error) {
	var files = make(map[string]map[string]string)

	if err := DownloadSources(b); err != nil {
		return files, fmt.Errorf("Could not download sources: %s", err)
	}
	srcdest := makepkg.MakepkgConf("SRCDEST")
	if srcdest == "" {
		srcdest = path.Join(b.ContainerPath, "srcdest")
	}
	b.Container.SetBindDir(srcdest, "/srcdest")
	pwd, err := os.Getwd()
	if err != nil {
		return files, err
	}
	b.Container.SetBindDir(pwd, "/startdir")
	if err := b.Container.Exec(makepkgCommand + makepkgArgs); err != nil {
		return files, err
	}
	files, err = utils.MoveProducts(b.Container)
	if err != nil {
		return files, err
	}
	return files, nil
}

// Init initializes the container
func (b *Builder) Init() error {
	// No need to initialize an initialized container
	if backend.CheckContainerExists(b.Path) {
		return nil
	}
	//
	ContainerPath, err := b.Backend.Setup()
	if err != nil {
		return err
	}
	// Point the container to the correct container
	b.ContainerPath = ContainerPath
	b.Container.SetPath(b.ContainerPath)
	//
	if err := b.Bootstrap.Init(b.ContainerPath); err != nil {
		return err
	}
	//
	if err := SetupConfig(b.ContainerPath, b.PacmanConf, b.MakepkgConf); err != nil {
		return err
	}
	// Create the initial files in the container
	if err := CreateFiles(b.ContainerPath, InitFileMap); err != nil {
		return err
	}
	// Generate locale files in the container
	if err := b.Container.Exec("locale-gen"); err != nil {
		return err
	}
	//
	if err := utils.SetupCacheDirs(b.Container, b.PacmanConf); err != nil {
		return err
	}
	// Upgrade packages in the container
	if err := b.Container.Exec("pacman -Syu --noconfirm base-devel"); err != nil {
		return err
	}
	return nil
}

// SetupChrootConfig -
func (b *Builder) SetupChrootConfig() error {
	if err := SetupGnupg(b.ContainerPath); err != nil {
		return err
	}
	return WithContainer(b.Container, ContainerActions{
		utils.SetupUser,
		utils.SetupMakepkgDirectories,
		utils.SetupMakepkg,
	})
}

// Fork - Sets up a snapshot from the root container
func (b *Builder) Fork(name string) error {
	if err := utils.SetupCacheDirs(b.Container, b.PacmanConf); err != nil {
		return err
	}
	if err := b.Container.Exec("pacman -Syu --noconfirm"); err != nil {
		return fmt.Errorf("Could not upgrade packages in container")
	}
	newContainerPath, err := b.Backend.AddSnapshot(name)
	if err != nil {
		return err
	}
	b.ContainerPath = newContainerPath
	b.Container.SetPath(newContainerPath)
	if err := b.SetupChrootConfig(); err != nil {
		return err
	}
	return nil
}

// Destroy - Removes a snapshot and defaults to root container
func (b *Builder) Destroy(name string) error {
	if err := b.Backend.RemoveSnapshot(name); err != nil {
		return err
	}
	b.ContainerPath = b.Path
	b.Container.SetPath(b.Path)
	return nil
}

// Update runs an update and copies over configs in case anything has changed
func (b *Builder) Update() error {
	return SetupConfig(b.ContainerPath, b.PacmanConf, b.MakepkgConf)
}

func NewBuilder() {
	panic("NewBuilder: Not implemented")
}
