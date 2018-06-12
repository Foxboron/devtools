package container

type Container interface {
	Exec(command string) error
	SetPath(path string)
	GetPath() string
	SetBindDir(src, dst string)
	SetBindRoDir(src, dst string)
}
