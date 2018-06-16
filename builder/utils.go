package builder

import (
	"github.com/foxboron/devtools/container"
	"github.com/foxboron/devtools/utils"
)

var (
	// Initial files (filename:contents)
	initFileMap = map[string][]byte{
		".arch-chroot":     []byte("v4"),
		"/etc/locale.conf": []byte("LANG=en_US.UTF-8"),
		"/etc/locale.gen":  []byte("en_US.UTF-8 UTF-8"),
	}
)

// CreateFiles takes a map over filenames and contents and creates
// those files in the container.
func CreateFiles(rootPath string, fileMap map[string][]byte) error {
	for filename, contents := range fileMap {
		if err := utils.CreateFile(rootPath, filename, contents, 0644); err != nil {
			return err
		}
	}
	return nil
}

// ContainerFunction can perform actions on a Container
type ContainerAction func(container.Container) error

// ContainerFunctions is a slice of ContainerFunction
type ContainerActions []ContainerAction

// WithContainer takes a slice of functions that takes a Container
// as an argument and applies them to the container within the Builder.
// If there is an error, the rest of the functions will not be called.
func WithContainer(container container.Container, cfs ContainerActions) error {
	c := container
	for _, cf := range cfs {
		if err := cf(c); err != nil {
			return err
		}
	}
	return nil
}
