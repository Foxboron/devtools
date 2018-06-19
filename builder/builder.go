package builder

import (
	"bufio"
	"errors"
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

// SetMirrorList reads the mirrorlist from pacman.conf and writes them to the
// container. Only includes servers used by the [core] repository.
func (b *Builder) SetMirrorList() error {
	file, err := os.Create(path.Join(b.ContainerPath, "etc", "pacman.d", "mirrorlist"))
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	defer w.Flush()

	pacmanconf, err := utils.GetPacmanConf(b.PacmanConf)
	if err != nil {
		return err
	}

	// We only want mirrors that are used by core so we don't
	// contaminate the container with random mirrors we don't want.
	for _, v := range pacmanconf.Repos {
		if v.Name != "core" {
			continue
		}
		for _, server := range v.Servers {
			fmt.Fprintln(w, "Server = "+server)
		}
		break
	}

	return nil
}

// SetupGnupg copies over the gnupg file from the host filesystem.
func (b *Builder) SetupGnupg() error {
	containerGnupg := path.Join(b.ContainerPath, "etc", "pacman.d", "gnupg")
	os.RemoveAll(containerGnupg)
	if _, err := utils.CopyDir(hostGnupgPath, containerGnupg); err != nil {
		return errors.New("Could not copy /etc/pacman.d/gnupg from host machine")
	}
	return nil
}

// SetupPacman copies over the pacman.conf and makepkg.conf from Builder settings
// and also scans the pacman.conf for any file inclusions we have to copy over
func (b *Builder) SetupPacman() error {
	if err := utils.CopyFile(b.PacmanConf, path.Join(b.ContainerPath, b.PacmanConf)); err != nil {
		return fmt.Errorf("Could not copy pacman.conf from %s", b.PacmanConf)
	}
	pacmanconf, err := utils.GetPacmanConf(b.PacmanConf)
	if err != nil {
		return err
	}
	for _, file := range pacmanconf.Include {
		if err := utils.CopyFile(file, path.Join(b.ContainerPath, file)); err != nil {
			return fmt.Errorf("Could not copy pacman.conf from %s", b.PacmanConf)
		}
	}
	if err := utils.CopyFile(b.MakepkgConf, path.Join(b.ContainerPath, b.MakepkgConf)); err != nil {
		return fmt.Errorf("Could not copy makepkg.conf from %s", b.MakepkgConf)
	}
	return nil
}

// SetupConfig sets up all the needed configurations for the container.
func (b *Builder) SetupConfig() error {
	if err := b.SetupGnupg(); err != nil {
		return err
	}
	if err := b.SetMirrorList(); err != nil {
		return err
	}
	if err := b.SetupPacman(); err != nil {
		return err
	}
	return nil
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
		b.ContainerPath = path.Join(b.Path, "root")
		b.Container.SetPath(b.ContainerPath)
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
	if err := b.SetupConfig(); err != nil {
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
	if err := b.SetupGnupg(); err != nil {
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
	return b.SetupConfig()
}

func NewBuilder() {
	panic("NewBuilder: Not implemented")
}
