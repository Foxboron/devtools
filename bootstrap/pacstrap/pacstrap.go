package pacstrap

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	alpm "github.com/Jguer/go-alpm"
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
	var cmd *exec.Cmd
	argArr := make([]string, 0)
	argArr = append(argArr, "sudo")
	argArr = append(argArr, "/usr/bin/pacstrap")
	argArr = append(argArr, "-"+p.Flags)
	if p.PacmanConf != "" {
		argArr = append(argArr, "-C "+p.PacmanConf)
	}
	// append root directory
	argArr = append(argArr, path)
	// Append packages we want
	for _, v := range p.Packages {
		argArr = append(argArr, v)
	}
	cmd = exec.Command(argArr[0], argArr[1:]...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Can't setup locale")
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
