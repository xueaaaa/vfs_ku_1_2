package session

import (
	"fmt"

	"vfs/src/internal/vfs"
)

type position struct {
	dirs  []*vfs.Directory
	names []string
}

func (p *position) current() *vfs.Directory {
	return p.dirs[len(p.dirs)-1]
}

func (p *position) up() {
	if len(p.names) == 0 {
		return
	}
	p.names = p.names[:len(p.names)-1]
	p.dirs = p.dirs[:len(p.dirs)-1]
}

func (p *position) down(name string) error {
	sub := findSubdir(p.current(), name)
	if sub == nil {
		return fmt.Errorf(msgNoSuchDir, name)
	}
	p.dirs = append(p.dirs, sub)
	p.names = append(p.names, name)
	return nil
}

func (p *position) step(part string) error {
	switch part {
	case "", CurrentDir:
		return nil
	case ParentDir:
		p.up()
		return nil
	}
	return p.down(part)
}

func findSubdir(dir *vfs.Directory, name string) *vfs.Directory {
	for i := range dir.Subdirs {
		if dir.Subdirs[i].Name == name {
			return &dir.Subdirs[i]
		}
	}
	return nil
}
