package bootstrap

import (
	"os"
)

type BootstrapType int

const (
	Archiso BootstrapType = iota
	Pacstrap
)

type Bootstrap interface {
	Init(path string) error
}

func DefaultBootstrap() BootstrapType {
	if _, err := os.Stat("/usr/bin/pacstrap"); os.IsNotExist(err) {
		return Archiso
	}
	return Pacstrap
}

func GetBootstrap(name string) BootstrapType {
	switch name {
	case "archiso":
		return Archiso
	case "pacstrap":
		return Pacstrap
	default:
		return DefaultBootstrap()
	}
}
