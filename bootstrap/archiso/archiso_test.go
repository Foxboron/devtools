package archiso

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"
)

func TestArchiso(t *testing.T) {
	// fmt.Println()

	dir, err := ioutil.TempDir("/var/tmp", "archiso")
	if err != nil {
		log.Fatal(err)
	}
	// defer os.RemoveAll(dir)
	fmt.Println(dir)

	archiso := newArchiso()
	archiso.TmpPath = "/var/tmp/archiso483364085"
	archiso.Path = dir
	archiso.Init()
}
