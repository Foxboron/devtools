package builder

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/foxboron/devtools/makepkg"
)

// DownloadSources drop priviledge to the non-sudo users and fetches the sources
func DownloadSources(builder *Builder) error {
	builddir, err := ioutil.TempDir("/var/tmp", "srcdir")
	if err != nil {
		return err
	}
	defer os.RemoveAll(builddir)
	// var c *exec.Cmd
	makepkgconf := path.Join(builder.ContainerPath, "etc", "makepkg.conf")
	var srcdest string
	if srcdest = makepkg.MakepkgConf("SRCDEST"); srcdest == "" {
		srcdest = path.Join(builder.ContainerPath, "srcdest")
	}
	cmdArgs := []string{"-u",
		os.Getenv("SUDO_USER"),
		"--preserve-env=GNUPGHOME",
		"env",
		fmt.Sprintf("SRCDEST=%s", srcdest),
		fmt.Sprintf("BUILDIR=%s", builddir),
		"makepkg",
		fmt.Sprintf("--config=%s", makepkgconf),
		"--verifysource",
		"-o"}
	var c *exec.Cmd
	c = exec.Command("/usr/bin/sudo", cmdArgs...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
