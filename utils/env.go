package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/foxboron/devtools/container"
	"github.com/foxboron/devtools/makepkg"
)

// SetupCacheDirs fetches all the CacheDirs from pacmans config and adds them
// as binds towards the container. The first cachedir is hopefully one we should
// be writing against. Everything else is added as read-only.
func SetupCacheDirs(container container.Container, PacmanConf string) error {
	pacmanconf, err := GetPacmanConf(PacmanConf)
	if err != nil {
		return err
	}
	container.SetBindDir(pacmanconf.CacheDir[0], pacmanconf.CacheDir[0])
	for _, cachedir := range pacmanconf.CacheDir[1:] {
		container.SetBindRoDir(cachedir, cachedir)
	}
	return nil
}

func CreateFile(containerPath, filename string, content []byte, mode os.FileMode) error {
	filePath := path.Join(containerPath, filename)
	path, _ := path.Split(filePath)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0760)
		if err != nil {
			return err
		}
	}
	err := ioutil.WriteFile(filePath, content, mode)
	if err != nil {
		return err
	}
	return nil
}

func SetupUser(container container.Container) error {
	uid := os.Getenv("SUDO_UID")
	// gid := os.Getenv("SUDO_GID")
	useradd := fmt.Sprintf("useradd -m -G wheel -d /build -u %s -s /bin/bash builduser", uid)
	err := container.Exec(useradd)
	if err != nil {
		return err
	}

	content := []byte("builduser ALL = NOPASSWD: /usr/bin/pacman\n")
	err = CreateFile(container.GetPath(), "/etc/sudoers.d/builduser-pacman", content, 0440)
	if err != nil {
		return err
	}

	return nil
}

// SetupMakepkgDirectories sets up the needed hook directories for makepkg
func SetupMakepkgDirectories(container container.Container) error {
	err := container.Exec("install -o builduser install -d {build,startdir,{pkg,srcpkg,src,log}dest}")
	if err != nil {
		return err
	}
	return nil
}

func SetupMakepkg(container container.Container) error {
	defaultValues := []string{
		"BUILDDIR=/build",
		"PKGDEST=/pkgdest",
		"SRCPKGDEST=/srcpkgdest",
		"SRCDEST=/srcdest",
		"LOGDEST=/logdest",
	}
	if makeflags := os.Getenv("MAKEFLAGS"); makeflags != "" {
		defaultValues = append(defaultValues, fmt.Sprintf("MAKEFLAGS=%s", makeflags))
	}

	if packager := os.Getenv("PACKAGER"); packager != "" {
		defaultValues = append(defaultValues, fmt.Sprintf("PACKAGER=%s", packager))
	}

	content := []byte(strings.Join(defaultValues, "\n") + "\n")
	err := CreateFile(container.GetPath(), "/build/.makepkg.conf", content, 0644)
	if err != nil {
		return err
	}
	return nil
}

// MoveProducts moves the created products from the container to
// the desired locations.
// Returns a map of {PKGDEST, LOGDEST, SRCPKGDEST} -> filename -> destination path
func MoveProducts(container container.Container) (map[string]map[string]string, error) {
	var products = make(map[string]map[string]string)

	// pkgdest
	pkgdest := makepkg.MakepkgConf("PKGDEST")
	if pkgdest == "" {
		pkgdest, _ = os.Getwd()
	}

	files, err := CopyDir(path.Join(container.GetPath(), "pkgdest"), pkgdest)
	if err != nil {
		return products, err
	}
	products["PKGDEST"] = files

	// logdest
	logdest := makepkg.MakepkgConf("LOGDEST")
	if logdest == "" {
		logdest, _ = os.Getwd()
	}
	files, err = CopyDir(path.Join(container.GetPath(), "logdest"), logdest)
	if err != nil {
		return products, err
	}
	products["LOGDEST"] = files

	// // srcpkgdest
	srcpkgdest := makepkg.MakepkgConf("SRCPKGDEST")
	if srcpkgdest == "" {
		srcpkgdest, _ = os.Getwd()
	}
	files, err = CopyDir(path.Join(container.GetPath(), "srcpkgdest"), srcpkgdest)
	if err != nil {
		return products, err
	}
	products["SRCPKGDEST"] = files

	return products, nil
}
