package pacstrap

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestPacstrap(t *testing.T) {
	fmt.Println()

	dir, err := ioutil.TempDir("/var/tmp", "pacstrap")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	fmt.Println(dir)

	pacstrap := newPacstrap()
	pacstrap.Path = dir
	pacstrap.init()
}
