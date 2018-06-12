package builder

import (
	"bufio"
	"fmt"
	"os"
	"path"

	"github.com/foxboron/devtools/backend"
	"github.com/foxboron/devtools/bootstrap"
	"github.com/foxboron/devtools/container"
	"github.com/foxboron/devtools/makepkg"
	"github.com/foxboron/devtools/utils"
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

// SetMirrorList reads the mirrorlist from pacman.conf and writes them to the container.
// Only includes servers used by the [core] repository.
func (b *Builder) SetMirrorList() error {
	file, err := os.Create(path.Join(b.ContainerPath, "etc", "pacman.d", "mirrorlist"))
	if err != nil {
		return err
	}
	defer file.Close()
	w := bufio.NewWriter(file)

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
	w.Flush()
	return nil
}

// SetupGnupg copies over the gnupg file from the host filesystem.
func (b *Builder) SetupGnupg() error {
	hostGnupg := path.Join("/etc", "pacman.d", "gnupg")
	containerGnupg := path.Join(b.ContainerPath, "etc", "pacman.d", "gnupg")
	os.RemoveAll(containerGnupg)
	err := utils.CopyDir(hostGnupg, containerGnupg)
	if err != nil {
		return fmt.Errorf("Couldnt copy /etc/pacman.d/gnupg from host machine")
	}
	return nil
}

// SetupPacman copies over the pacman.conf and makepkg.conf from Builder settings
// and also scans the pacman.conf for any file inclusions we have to copy over
func (b *Builder) SetupPacman() error {
	err := utils.CopyFile(b.PacmanConf, path.Join(b.ContainerPath, b.PacmanConf))
	if err != nil {
		return fmt.Errorf("Couldnt copy pacman.conf from %s", b.PacmanConf)
	}

	pacmanconf, err := utils.GetPacmanConf(b.PacmanConf)
	if err != nil {
		return err
	}

	for _, file := range pacmanconf.Include {
		err := utils.CopyFile(file, path.Join(b.ContainerPath, file))
		if err != nil {
			return fmt.Errorf("Couldnt copy pacman.conf from %s", b.PacmanConf)
		}
	}

	err = utils.CopyFile(b.MakepkgConf, path.Join(b.ContainerPath, b.MakepkgConf))
	if err != nil {
		return fmt.Errorf("Couldnt copy makepkg.conf from %s", b.MakepkgConf)
	}
	return nil
}

// SetupConfig sets up all the needed configurations for the container.
func (b *Builder) SetupConfig() error {
	err := b.SetupGnupg()
	if err != nil {
		return err
	}

	err = b.SetMirrorList()
	if err != nil {
		return err
	}

	err = b.SetupPacman()
	if err != nil {
		return err
	}

	return nil
}

func (b *Builder) Build() error {
	err := DownloadSources(b)
	if err != nil {
		//return fmt.Errorf("Can't download sources")
	}
	var srcdest string
	if srcdest = makepkg.MakepkgConf("SRCDEST"); srcdest == "" {
		srcdest = path.Join(b.ContainerPath, "srcdest")
	}
	b.Container.SetBindDir(srcdest, "/srcdest")
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	makepkgArgs := "syncdeps --noconfirm --log --holdver --skipinteg"
	b.Container.SetBindDir(pwd, "/startdir")
	err = b.Container.Exec("sudo --preserve-env=SOURCE_DATE_EPOCH -iu builduser bash -c 'cd /startdir; makepkg \"$@\"' -bash " + makepkgArgs)
	if err != nil {
		return err
	}
	err = utils.MoveProducts(b.Container)
	if err != nil {
		return err
	}
	return nil
}

// Init initializes the container
func (b *Builder) Init() error {
	// No need to initialize an initialized container

	exists := backend.CheckContainerExists(b.Path)
	if exists {
		b.ContainerPath = path.Join(b.Path, "root")
		b.Container.SetPath(b.ContainerPath)
		return nil
	}

	ContainerPath, err := b.Backend.Setup()
	if err != nil {
		return err
	}

	// Point the container to the correct container
	b.ContainerPath = ContainerPath
	b.Container.SetPath(b.ContainerPath)

	err = b.Bootstrap.Init(b.ContainerPath)
	if err != nil {
		return err
	}

	err = b.SetupConfig()
	if err != nil {
		return err
	}

	err = utils.CreateFile(b.ContainerPath, ".arch-chroot", []byte("v4"), 0644)
	if err != nil {
		return err
	}
	err = utils.CreateFile(b.ContainerPath, "/etc/locale.conf", []byte("LANG=en_US.UTF-8"), 0644)
	if err != nil {
		return err
	}
	err = utils.CreateFile(b.ContainerPath, "/etc/locale.gen", []byte("en_US.UTF-8 UTF-8"), 0644)
	if err != nil {
		return err
	}
	err = b.Container.Exec("locale-gen")
	if err != nil {
		return err
	}

	err = utils.SetupCacheDirs(b.Container, b.PacmanConf)
	if err != nil {
		return err
	}

	err = b.Container.Exec("pacman -Syu --noconfirm base-devel")
	if err != nil {
		return err
	}
	return nil
}

func (b *Builder) SetupChrootConfig() error {

	err := b.SetupGnupg()
	if err != nil {
		return err
	}

	err = utils.SetupUser(b.Container)
	if err != nil {
		return err
	}

	err = utils.SetupMakepkgDirectories(b.Container)
	if err != nil {
		return err
	}

	err = utils.SetupMakepkg(b.Container)
	if err != nil {
		return err
	}
	return nil
}

// Fork - Sets up a snapshot from the root container
func (b *Builder) Fork(name string) error {
	err := utils.SetupCacheDirs(b.Container, b.PacmanConf)
	if err != nil {
		return err
	}

	err = b.Container.Exec("pacman -Syu --noconfirm")
	if err != nil {
		return fmt.Errorf("Couldn't update container")
	}
	newContainerPath, err := b.Backend.AddSnapshot(name)
	if err != nil {
		return err
	}
	b.ContainerPath = newContainerPath
	b.Container.SetPath(newContainerPath)
	err = b.SetupChrootConfig()
	if err != nil {
		return err
	}
	return nil
}

// Destroy - Removes a snapshot and defaults to root container
func (b *Builder) Destroy(name string) error {
	err := b.Backend.RemoveSnapshot(name)
	if err != nil {
		return err
	}
	b.ContainerPath = path.Join(b.Path, "root")
	b.Container.SetPath(path.Join(b.Path, "root"))
	return nil
}

// Update runs an update and copies over configs incase anything has changed
func (b *Builder) Update() error {
	err := b.SetupConfig()
	if err != nil {
		return err
	}
	return nil
}

func NewBuilder() {
}
