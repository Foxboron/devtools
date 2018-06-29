package pacstrap

import (
	"fmt"
	"log"
	"os"
	"errors"

	alpm "github.com/Jguer/go-alpm"
	"github.com/xyproto/argbuilder"
)

// Defaults
var (
	PacmanConf string = "/etc/pacman.conf"
)

type Pacstrap struct {
	Name       string
	Path       string
	Cachedirs  []string
	PacmanConf string
	Flags      string // GMcd
	Packages   []string
}

func (p *Pacstrap) Init(path string) error {
	ab := argbuilder.New("sudo", "/usr/bin/pacstrap", "-" + p.Flags)
	if p.PacmanConf != "" {
		ab.Add("-C " + p.PacmanConf)
	}
	// append root directory
	ab.Add(path)
	// Append packages we want
	ab.AddStrings(p.Packages)
	if ab.Run() != nil {
		return errors.New("Can't setup locale")
	}
	return nil
}

func NewPacstrap(PacmanConf string) *Pacstrap {
	b, err := os.Open(PacmanConf)
	if err != nil {
		fmt.Print(err)
	}
	defer b.Close()

	pacmanconf, err := alpm.ParseConfig(b)
	if err != nil {
		log.Fatal(err)
	}

	return &Pacstrap{
		Name:       "",
		Path:       "",
		Flags:      "GMcd",
		PacmanConf: PacmanConf,
		Cachedirs:  pacmanconf.CacheDir,
		Packages:   []string{"base-devel"},
	}

}
