package fs

import (
	"fmt"
	"os"
	"path"
)

type Filesystem struct {
	Name string
	Path string
}

func (fs *Filesystem) setup() error {
	err := os.MkdirAll(path.Join(fs.Path, "init"), 0700)
	if err != nil {
		fmt.Println("Failed to create directory")
	}
	return nil
}

func (fs *Filesystem) addSnapshot() error {
	return nil
}

func (fs *Filesystem) removeSnapshot() {
	err := os.Remove(path.Join(fs.Path, "init"))
	if err != nil {
		fmt.Println("Failed to remove directory")
	}
}

func (fs *Filesystem) destroy() {
	err := os.Remove(path.Join(fs.Path, "init"))
	if err != nil {
		fmt.Println("Failed to remove directory")
	}
}
