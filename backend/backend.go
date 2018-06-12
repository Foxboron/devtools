package backend

import (
	"io/ioutil"
	"log"
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

func GetBackendFromContainer(containerPath string) BackendFilesystem {
	fsFile := path.Join(containerPath, "root", ".arch-chroot-fs")
	if _, err := os.Stat(fsFile); os.IsNotExist(err) {
		return -1
	}
	buf, err := ioutil.ReadFile(fsFile)
	if err != nil {
		log.Fatal(err)
	}
	return GetBackend(string(buf))
}

func CheckContainerExists(containerPath string) bool {
	fsFile := path.Join(containerPath, "root", ".arch-chroot-fs")
	if _, err := os.Stat(fsFile); os.IsNotExist(err) {
		return false
	}
	return true
}
