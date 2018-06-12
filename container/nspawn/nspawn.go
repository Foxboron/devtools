package nspawn

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
)

var (
	stdout io.Writer
	stderr io.Writer
	stdin  io.Reader
)

func init() {
	stdout = os.Stdout
	stderr = os.Stderr
	stdin = nil
}

type Nspawn struct {
	// Env        *environment.Environment
	Path       string
	BindDirs   map[string]string
	BindRoDirs map[string]string
	Flags      []string
}

func (n *Nspawn) Exec(command string) error {
	var c *exec.Cmd
	cmdArgs := []string{}
	cmdArgs = append(cmdArgs, "-D", n.Path)
	cmdArgs = append(cmdArgs, n.Flags...)
	cmdArgs = append(cmdArgs, n.FormatBind()...)
	cmdArgs = append(cmdArgs, "/bin/sh", "-c")
	cmdArgs = append(cmdArgs, command)
	c = exec.Command("systemd-nspawn", cmdArgs...)
	c.Stdout = stdout
	c.Stderr = stderr
	c.Stdin = stdin
	return c.Run()
}

// ExecScript - Excecute a script
func (n *Nspawn) ExecScript(script []string) (string, error) {
	return "", nil
}

func (n *Nspawn) InsertScript(name string, script []string) error {
	return nil
}

func (n *Nspawn) setMachineId() {
	if n.Path == "" {
		return
	}
	if _, err := os.Stat(path.Join(n.Path, "etc", "machine-id")); os.IsNotExist(err) {
		var c *exec.Cmd
		c = exec.Command("systemd-machine-id-setup", fmt.Sprintf("--root=%s", n.Path))
		c.Run()
	}
}

func (n *Nspawn) SetPath(path string) {
	n.Path = path
	n.setMachineId()
}

func (n *Nspawn) GetPath() string {
	return n.Path
}

func (n *Nspawn) SetBindDir(src, dst string) {
	n.BindDirs[src] = dst
}

func (n *Nspawn) SetBindRoDir(src, dst string) {
	n.BindRoDirs[src] = dst
}

func (n *Nspawn) FormatBind() []string {
	var bindList []string
	for src, dest := range n.BindDirs {
		if src == dest {
			bindList = append(bindList, fmt.Sprintf("--bind=%s", src))
		} else {
			bindList = append(bindList, fmt.Sprintf("--bind=%s:%s", src, dest))
		}
	}
	for src, dest := range n.BindRoDirs {
		if src == dest {
			bindList = append(bindList, fmt.Sprintf("--bind-ro=%s", src))
		} else {
			bindList = append(bindList, fmt.Sprintf("--bind-ro=%s:%s", src, dest))
		}
	}
	return bindList
}

func NewNspawn(path string) *Nspawn {
	n := &Nspawn{
		Path: path,
		Flags: []string{
			"-q",
			"--as-pid2",
			"--register=no",
		},
		BindDirs:   make(map[string]string),
		BindRoDirs: make(map[string]string),
	}
	n.setMachineId()
	return n
}
