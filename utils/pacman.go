package utils

import (
	"fmt"
	"os"

	alpm "github.com/Jguer/go-alpm"
)

func GetPacmanConf(path string) (alpm.PacmanConfig, error) {
	buf, err := os.Open(path)
	if err != nil {
		return alpm.PacmanConfig{}, fmt.Errorf("Couldnt read pacman.conf from %s", path)
	}
	defer buf.Close()

	pacmanconf, err := alpm.ParseConfig(buf)
	if err != nil {
		return alpm.PacmanConfig{}, fmt.Errorf("Couldnt parse pacman.conf from %s", path)
	}
	return pacmanconf, nil
}
