package archiso

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	alpm "github.com/Jguer/go-alpm"
	"github.com/foxboron/devtools/bootstrap"
	"github.com/foxboron/devtools/utils"
)

var (
	PacmanConf  string = "/etc/pacman.conf"
	IsoCacheDir        = "/var/cache/devtools"
)

type Archiso struct {
	Mirror  string
	ISOName string
	Path    string

	// This is where we store the ISO
	TmpPath string
}

func (a *Archiso) DownloadISO() (string, error) {
	if a.TmpPath == "" {
		err := os.MkdirAll(IsoCacheDir, 0755)
		if err != nil {
			log.Fatal(err)
		}
		a.TmpPath = IsoCacheDir
	}
	isoPath := path.Join(a.TmpPath, a.ISOName)
	if _, err := os.Stat(isoPath); os.IsNotExist(err) {
		utils.Msg2f("Downloading %s...", a.ISOName)
		utils.DownloadFile(isoPath, a.Mirror+a.ISOName)
		utils.DownloadFile(isoPath+".sig", a.Mirror+a.ISOName+".sig")
		err := utils.VerifySignature(path.Join(a.TmpPath, a.ISOName), path.Join(a.TmpPath, a.ISOName+".sig"))
		if err != nil {
			utils.Error("Signature check failed!")
			os.Exit(1)
		}
	}
	return isoPath, nil
}

func (a *Archiso) Init(dst string) error {
	isoPath, err := a.DownloadISO()
	if err != nil {
		return err
	}
	if _, err := os.Stat(path.Join(dst, ".arch-chroot")); os.IsNotExist(err) {
		fmt.Println(isoPath)
		fmt.Println(dst)
		err = utils.Untar(isoPath, dst)
		if err != nil {
			return fmt.Errorf("Could not untar ISO: %s", err)
		}
	}
	return nil
}

func GetArchIsoName() string {
	name := "archlinux-bootstrap-%s.01-%s.tar.gz"
	t := time.Now()
	return fmt.Sprintf(name, t.Format("2006.01"), "x86_64")
}

func NewArchiso(PacmanConf string) bootstrap.Bootstrap {
	b, err := os.Open(PacmanConf)
	if err != nil {
		fmt.Print(err)
	}

	pacmanconf, err := alpm.ParseConfig(b)
	if err != nil {
		log.Fatal(err)
	}

	// Get the first mirror to core
	var mirror string
	for _, v := range pacmanconf.Repos {
		if v.Name != "core" {
			continue
		}
		s := strings.Split(v.Servers[0], "$repo")
		mirror = s[0]
		break
	}

	return &Archiso{
		Mirror:  mirror + "iso/latest/",
		ISOName: GetArchIsoName(),
		TmpPath: "",
		Path:    "",
	}
}
