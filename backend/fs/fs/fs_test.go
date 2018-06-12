package fs

import (
	"testing"
)

var fs Filesystem

func init() {
	fs = Filesystem{Name: "Test", Path: "./test"}
}

func TestInit(t *testing.T) {
	fs.init()
}

func TestRemove(t *testing.T) {
	fs.remove()
}
