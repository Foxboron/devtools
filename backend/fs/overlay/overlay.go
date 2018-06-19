package overlay

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"syscall"

	"github.com/foxboron/devtools/backend"
	"github.com/foxboron/devtools/utils"
)

type Overlay struct {
	Name        string
	RootPath    string
	CurrentPath string
	Directories map[string]string
}

func (o *Overlay) Setup() (string, error) {
	err := os.MkdirAll(o.RootPath, 0755)
	if err != nil {
		return "", fmt.Errorf("Failed to setup overlay backend")
	}
	fileInfo := path.Join(o.RootPath, ".arch-chroot-fs")
	err = ioutil.WriteFile(fileInfo, []byte("overlay"), 0644)
	if err != nil {
		return "", fmt.Errorf("Failed to write filesystem file")
	}
	return o.RootPath, nil
}

func (o *Overlay) AddSnapshot(name string) (string, error) {
	o.Directories["root"] = o.RootPath
	directory, _ := path.Split(o.RootPath)
	o.Directories["build"] = path.Join(directory, name)
	o.Directories["upperdir"] = path.Join(directory, name+"_upperdir")
	o.Directories["workdir"] = path.Join(directory, name+"_workdir")
	for _, dirPath := range o.Directories {
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			return "", fmt.Errorf("Failed to setup overlay backend: %s", err)
		}
	}
	flags := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s",
		o.Directories["root"],
		o.Directories["upperdir"],
		o.Directories["workdir"],
	)
	err := syscall.Mount(o.Directories["root"], o.Directories["build"], "overlay", 0, flags)
	if err != nil {
		return "", fmt.Errorf("Failed to mount with overlay: %s", err)
	}
	return o.Directories["build"], nil
}

func (o *Overlay) RemoveSnapshot(name string) error {
	err := syscall.Unmount(o.Directories["build"], 0)
	if err != nil {
		log.Fatal(err)
	}
	for key, dirPath := range o.Directories {
		if key == "root" {
			continue
		}
		err := os.RemoveAll(dirPath)
		if err != nil {
			return fmt.Errorf("Failed to cleanup overlay")
		}
	}
	return nil
}

func (o *Overlay) Destroy() error {
	err := syscall.Unmount(o.Directories["build"], 0)
	if err != nil {
		utils.Warning("Filesystem has been unmounted allready")
	}
	for _, dirPath := range o.Directories {
		err := os.RemoveAll(dirPath)
		if err != nil {
			return fmt.Errorf("Failed to cleanup overlay")
		}
	}
	return nil
}

func (o *Overlay) GetPath() string {
	return o.CurrentPath
}

func NewOverlay(path string) backend.Backend {
	return &Overlay{
		RootPath:    path,
		CurrentPath: path,
		Directories: make(map[string]string),
	}
}
