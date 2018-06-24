package builder

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/foxboron/devtools/utils"
)

// SetupPacman copies over the pacman.conf and makepkg.conf from Builder settings
// and also scans the pacman.conf for any file inclusions we have to copy over
func SetupPacman(containerPath, PacmanConf, MakepkgConf string) error {
	if err := utils.CopyFile(PacmanConf, path.Join(containerPath, PacmanConf)); err != nil {
		return fmt.Errorf("Could not copy pacman.conf from %s", PacmanConf)
	}
	pacmanconf, err := utils.GetPacmanConf(PacmanConf)
	if err != nil {
		return err
	}
	for _, file := range pacmanconf.Include {
		if err := utils.CopyFile(file, path.Join(containerPath, file)); err != nil {
			return fmt.Errorf("Could not copy pacman.conf from %s", PacmanConf)
		}
	}
	if err := utils.CopyFile(MakepkgConf, path.Join(containerPath, MakepkgConf)); err != nil {
		return fmt.Errorf("Could not copy makepkg.conf from %s", MakepkgConf)
	}
	return nil
}

// SetMirrorList reads the mirrorlist from pacman.conf and writes them to the
// container. Only includes servers used by the [core] repository.
func SetMirrorList(containerPath string, PacmanConf string) error {
	file, err := os.Create(path.Join(containerPath, "etc", "pacman.d", "mirrorlist"))
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	defer w.Flush()

	pacmanconf, err := utils.GetPacmanConf(PacmanConf)
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
func SetupGnupg(containerPath string) error {
	containerGnupg := path.Join(containerPath, "etc", "pacman.d", "gnupg")
	os.RemoveAll(containerGnupg)
	if _, err := utils.CopyDir(hostGnupgPath, containerGnupg); err != nil {
		return errors.New("Could not copy /etc/pacman.d/gnupg from host machine")
	}
	return nil
}

// SetupConfig sets up all the needed configurations for the container.
func SetupConfig(containerPath, PacmanConf, MakepkgConf string) error {
	if err := SetupGnupg(containerPath); err != nil {
		return err
	}
	if err := SetMirrorList(containerPath, PacmanConf); err != nil {
		return err
	}
	return SetupPacman(containerPath, PacmanConf, MakepkgConf)
}
