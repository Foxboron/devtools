package backend

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

type BackendFilesystem int

const (
	Fs BackendFilesystem = iota
	Btrfs
	Overlay
	Rsync
)

type Backend interface {
	Setup() (string, error)
	AddSnapshot(name string) (string, error)
	RemoveSnapshot(name string) error
	Destroy() error
	GetPath() string
}

// GetBackend - get a backend with string name
func GetBackend(name string) BackendFilesystem {
	switch name {
	case "overlay":
		return Overlay
	// case "rsync":
	// 	return &rsync.Rsync{}, nil
	default:
		return -1
	}
}

func GetBackendFromContainer(containerPath string) (BackendFilesystem, error) {
	fsFile := path.Join(containerPath, ".arch-chroot-fs")
	if _, err := os.Stat(fsFile); os.IsNotExist(err) {
		return -1, fmt.Errorf("Not a container")
	}
	buf, err := ioutil.ReadFile(fsFile)
	if err != nil {
		return -1, fmt.Errorf("Not a container")
	}
	return GetBackend(string(buf)), nil
}

func CheckContainerExists(containerPath string) bool {
	fsFile := path.Join(containerPath, "root", ".arch-chroot-fs")
	if _, err := os.Stat(fsFile); os.IsNotExist(err) {
		return false
	}
	return true
}
