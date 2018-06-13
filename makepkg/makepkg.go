package makepkg

import (
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"strings"
)

var (
	env              = make(map[string]string)
	makepkgCache     = make(map[string]string)
	userMakepkgCache = make(map[string]string)
)

func GetMakepkgConfFile() string {
	// var content string
	return ""
}

// MakepkgConf checks first the env for variables, then defers to the user makepkg.conf before the system makepkg.conf
// is checked for the value
func MakepkgConf(key string) string {
	if val := os.Getenv(key); key != "" {
		return val
	}
	if val, ok := userMakepkgCache[key]; ok {
		return val
	}

	if val, ok := makepkgCache[key]; ok {
		return val
	}
	return ""
}

func ParseMakepkgConf(path string) map[string]string {
	var config = make(map[string]string)

	file, _ := ioutil.ReadFile(path)
	lines := strings.Split(string(file), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		split := strings.SplitN(line, "=", 2)
		if len(split) != 2 {
			continue
		}

		switch split[0] {
		case "SRCDEST":
			config[split[0]] = split[1]
		case "SRCPKGDEST":
			config[split[0]] = split[1]
		case "PKGDEST":
			config[split[0]] = split[1]
		case "LOGDEST":
			config[split[0]] = split[1]
		case "MAKEFLAGS":
			config[split[0]] = split[1]
		case "PACKAGER":
			config[split[0]] = split[1]
		}
	}
	return config
}

func InitMakepkgConf() {
	// TODO: Should probably not assume Arch Linux
	makepkgCache = ParseMakepkgConf("/etc/makepkg.conf")

	// user
	var configPath string
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		configPath = xdg
	} else {
		uid := os.Getenv("SUDO_UID")
		if uid != "" {
			user, err := user.LookupId(uid)
			if err != nil {
				log.Fatal(err)
			}
			configPath = user.HomeDir
		} else {
			user, err := user.Current()
			if err != nil {
				log.Fatal(err)
			}
			configPath = user.HomeDir
		}
	}
	makepkgConfPath := path.Join(configPath, "config", "pacman", "makepkg.conf")
	if _, err := os.Stat(makepkgConfPath); err == nil {
		userMakepkgCache = ParseMakepkgConf(makepkgConfPath)
		return
	}

	makepkgConfPath = path.Join(configPath, ".makepkg.conf")
	if _, err := os.Stat(makepkgConfPath); err == nil {
		userMakepkgCache = ParseMakepkgConf(makepkgConfPath)
		return
	}
}

func init() {
	InitMakepkgConf()
}
