package session

import (
	"fmt"
	"path"
	"strings"

	"vfs/src/internal/vfs"
)

const (
	PathSep    = "/"
	CurrentDir = "."
	ParentDir  = ".."

	msgNoSuchDir  = "no such directory: %s"
	msgNoSuchFile = "file not found: %s"
	msgBadPath    = "invalid path: %s"
)

type Session struct {
	Root    *vfs.Directory
	Current *vfs.Directory
	Path    []string
}

func NewSession(root *vfs.Directory) *Session {
	return &Session{Root: root, Current: root, Path: []string{}}
}

func (s *Session) PathString() string {
	return PathSep + strings.Join(s.Path, PathSep)
}

func (s *Session) Cd(target string) error {
	pos, err := s.resolve(target)
	if err != nil {
		return err
	}
	s.Current, s.Path = pos.current(), pos.names
	return nil
}

func (s *Session) ResolveDir(target string) (*vfs.Directory, error) {
	pos, err := s.resolve(target)
	if err != nil {
		return nil, err
	}
	return pos.current(), nil
}

func (s *Session) FindFile(target string) (*vfs.File, error) {
	parent, name, err := s.locate(target)
	if err != nil {
		return nil, err
	}
	dir := parent.current()
	for i := range dir.Files {
		if dir.Files[i].Name == name {
			return &dir.Files[i], nil
		}
	}
	return nil, fmt.Errorf(msgNoSuchFile, name)
}

func (s *Session) locate(target string) (*position, string, error) {
	dirPart, name := path.Split(target)
	if name == "" {
		return nil, "", fmt.Errorf(msgBadPath, target)
	}
	parent, err := s.resolve(dirPart)
	if err != nil {
		return nil, "", err
	}
	return parent, name, nil
}

func (s *Session) resolve(target string) (*position, error) {
	pos, err := s.start(target)
	if err != nil {
		return nil, err
	}
	for _, part := range strings.Split(target, PathSep) {
		if err := pos.step(part); err != nil {
			return nil, err
		}
	}
	return pos, nil
}

func (s *Session) start(target string) (*position, error) {
	pos := &position{dirs: []*vfs.Directory{s.Root}}
	if strings.HasPrefix(target, PathSep) {
		return pos, nil
	}
	for _, name := range s.Path {
		if err := pos.down(name); err != nil {
			return nil, err
		}
	}
	return pos, nil
}
